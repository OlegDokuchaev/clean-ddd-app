package otp_store

import (
	"context"
	"fmt"
	"strconv"

	customerAplication "customer/internal/application/customer"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type OtpStoreImpl struct {
	cfg *Config
	rdb *redis.Client
}

func New(cfg *Config, rdb *redis.Client) *OtpStoreImpl {
	return &OtpStoreImpl{cfg: cfg, rdb: rdb}
}

func (s *OtpStoreImpl) key(challengeID string) string {
	return s.cfg.KeyPrefix + challengeID
}

func (s *OtpStoreImpl) Issue(ctx context.Context, challengeID string, consumerID uuid.UUID, code string, policy customerAplication.OtpPolicy) error {
	key := s.key(challengeID)

	pipe := s.rdb.TxPipeline()
	pipe.HSet(ctx, key, map[string]any{
		"code":          code,
		"consumer_id":   consumerID.String(),
		"attempts_left": policy.MaxAttempts,
	})
	pipe.Expire(ctx, key, policy.TTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrOtpStoreUnavailable, err)
	}

	return nil
}

func (s *OtpStoreImpl) VerifyAndConsume(
	ctx context.Context,
	challengeID string,
	code string,
) (ok bool, attemptsLeft int, expired bool, consumerID uuid.UUID, err error) {
	key := s.key(challengeID)

	var (
		outOK          bool
		outAttempts    int
		outExpired     bool
		outConsumerStr string
	)

	wErr := s.rdb.Watch(ctx, func(tx *redis.Tx) error {
		h, err := tx.HGetAll(ctx, key).Result()
		if err != nil {
			return ErrOtpStoreUnavailable
		}
		if len(h) == 0 {
			outOK, outAttempts, outExpired, outConsumerStr = false, 0, true, ""
			return nil
		}

		storedCode, okCode := h["code"]
		attStr, okAttempts := h["attempts_left"]
		outConsumerStr = h["consumer_id"]

		if !okCode || !okAttempts {
			return ErrOtpStoreCorrupted
		}
		curAttempts, convErr := strconv.Atoi(attStr)
		if convErr != nil {
			return ErrOtpStoreCorrupted
		}

		if curAttempts <= 0 {
			outOK, outAttempts, outExpired = false, 0, false
			return nil
		}

		if storedCode == code {
			_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
				p.Del(ctx, key)
				return nil
			})
			if err != nil {
				return ErrOtpStoreUnavailable
			}
			outOK, outAttempts, outExpired = true, curAttempts, false
			return nil
		}

		newAttempts := curAttempts - 1
		if newAttempts < 0 {
			newAttempts = 0
		}
		_, err = tx.TxPipelined(ctx, func(p redis.Pipeliner) error {
			p.HSet(ctx, key, "attempts_left", newAttempts)
			return nil
		})
		if err != nil {
			return ErrOtpStoreUnavailable
		}

		outOK, outAttempts, outExpired = false, newAttempts, false
		return nil
	}, key)
	if wErr != nil {
		return false, 0, false, uuid.Nil, wErr
	}

	consumerID, err = uuid.Parse(outConsumerStr)
	if err != nil {
		return false, 0, false, uuid.Nil, err
	}

	return outOK, outAttempts, outExpired, consumerID, nil
}

func (s *OtpStoreImpl) Invalidate(ctx context.Context, challengeID string) error {
	key := s.key(challengeID)
	if err := s.rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("%w: %v", ErrOtpStoreUnavailable, err)
	}
	return nil
}

var _ customerAplication.OtpStore = (*OtpStoreImpl)(nil)

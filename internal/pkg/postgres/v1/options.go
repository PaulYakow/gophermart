package v1

import "time"

type Option func(*Postgre)

func MaxPoolSize(size int) Option {
	return func(p *Postgre) {
		p.maxPollSize = size
	}
}

func MaxConnIdleTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnIdleTime = duration
	}
}

func MaxConnLifeTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnLifeTime = duration
	}
}

func ConnAttempts(attempts int) Option {
	return func(p *Postgre) {
		p.connAttempts = attempts
	}
}

func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgre) {
		p.connTimeout = timeout
	}
}

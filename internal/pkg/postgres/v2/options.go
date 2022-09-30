package v2

import "time"

type Option func(*Postgre)

func MaxOpenConn(size int) Option {
	return func(p *Postgre) {
		p.maxOpenConn = size
	}
}

func MaxIdleConn(size int) Option {
	return func(p *Postgre) {
		p.maxIdleConn = size
	}
}

func MaxConnIdleTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnIdleTime = duration
	}
}

func MaxConnLifeTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnLifetime = duration
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

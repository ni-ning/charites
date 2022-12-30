package bootstrap

import "log"

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupLogger()
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("init.setupDBEngine err: %v", err)
	}

	err = setupRedis()
	if err != nil {
		log.Fatalf("init.setupRedis err: %v", err)
	}

	err = setupSnowflake("", 1)
	if err != nil {
		log.Fatalf("init.setupSnowflake err: %v", err)
	}

	err = setupRPClient()
	if err != nil {
		log.Fatalf("init.setupRPClient err: %v", err)
	}

	err = setupRocketMQ()
	if err != nil {
		log.Fatalf("setupRocketMQ err: %v", err)
	}

}

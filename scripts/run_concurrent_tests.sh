#!/usr/bin/env bash
run_test_by_name () {

	cd ../internal/broker
	for i in range $(seq 1 50)
	do
		go test -v -run "$1"
	done

}

run_test_by_name TestConcurrentPublishOnOneSubjectShouldNotFail
run_test_by_name TestConcurrentSubscribesShouldNotFail
run_test_by_name TestConcurrentPublishOnOneSubjectShouldNotFail
run_test_by_name TestConcurrentPublishShouldNotFail



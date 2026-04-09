# Performance & Scalability Summary (k6)

**Environment:** Docker Compose (single-node)
**Scope:** Read-only endpoints
**Design Principles:** Bounded concurrency, backpressure, graceful degradation

## User Capability

```text
█ TOTAL RESULTS

    checks_total.......: 173384 962.687169/s
    checks_succeeded...: 99.99% 173368 out of 173384
    checks_failed......: 0.00%  16 out of 173384

    ✓ GET /api/user?page=1&limit=10 success
    ✓ GET /api/user/41 success
    ✗ GET /api/user/active?page=1&limit=10 success
      ↳  99% — ✓ 43341 / ✗ 5
    ✗ GET /api/user/trashed?page=1&limit=10 success
      ↳  99% — ✓ 43335 / ✗ 11

    HTTP
    http_req_duration..............: avg=3.72ms   min=513.71µs med=911.1µs  max=15.01s p(90)=3.02ms   p(95)=6.74ms
      { expected_response:true }...: avg=2.33ms   min=513.71µs med=911.06µs max=5.16s  p(90)=3.01ms   p(95)=6.73ms
    http_req_failed................: 0.00%  16 out of 173384
    http_reqs......................: 173384 962.687169/s

    EXECUTION
    dropped_iterations.............: 152    0.843956/s
    iteration_duration.............: avg=116.15ms min=103.08ms med=104.7ms  max=15.17s p(90)=114.41ms p(95)=130.21ms
    iterations.....................: 43346  240.671792/s
    vus............................: 64     min=5            max=137
    vus_max........................: 150    min=100          max=150

    NETWORK
    data_received..................: 406 MB 2.3 MB/s
    data_sent......................: 45 MB  252 kB/s




running (3m00.1s), 000/150 VUs, 43346 complete and 0 interrupted iterations
scalability_test ✓ [======================================] 000/150 VUs  3m0s  599.99 iters/s
```

## User Load Test

```text
 █ THRESHOLDS

    http_req_duration
    ✗ 'p(95)<200' p(95)=385.69ms

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%


  █ TOTAL RESULTS

    checks_total.......: 460948  3816.062108/s
    checks_succeeded...: 100.00% 460948 out of 460948
    checks_failed......: 0.00%   0 out of 460948

    ✓ GET /api/user?page=1&limit=10 success
    ✓ GET /api/user/41 success
    ✓ GET /api/user/active?page=1&limit=10 success
    ✓ GET /api/user/trashed?page=1&limit=10 success

    HTTP
    http_req_duration..............: avg=234.92ms min=43.03ms  med=210.27ms max=1.38s p(90)=331.31ms p(95)=385.69ms
      { expected_response:true }...: avg=234.92ms min=43.03ms  med=210.27ms max=1.38s p(90)=331.31ms p(95)=385.69ms
    http_req_failed................: 0.00%  0 out of 460948
    http_reqs......................: 460948 3816.062108/s

    EXECUTION
    iteration_duration.............: avg=1.04s    min=694.65ms med=966.23ms max=3.16s p(90)=1.34s    p(95)=1.5s
    iterations.....................: 115237 954.015527/s
    vus............................: 51     min=51          max=1000
    vus_max........................: 1000   min=1000        max=1000

    NETWORK
    data_received..................: 1.1 GB 8.9 MB/s
    data_sent......................: 121 MB 998 kB/s




running (2m00.8s), 0000/1000 VUs, 115237 complete and 0 interrupted iterations
load_test ✓ [======================================] 1000 VUs  2m0s
ERRO[0121] thresholds on metrics 'http_req_duration' have been crossed
```

## User Rampling Test

```text
█ TOTAL RESULTS

    checks_total.......: 431748 3568.709106/s
    checks_succeeded...: 99.99% 431718 out of 431748
    checks_failed......: 0.00%  30 out of 431748

    ✗ GET /api/user?page=1&limit=10 success
      ↳  99% — ✓ 107907 / ✗ 30
    ✓ GET /api/user/41 success
    ✓ GET /api/user/active?page=1&limit=10 success
    ✓ GET /api/user/trashed?page=1&limit=10 success

    HTTP
    http_req_duration..............: avg=163.67ms min=570.64µs med=134.49ms max=15.07s p(90)=348.2ms  p(95)=408.28ms
      { expected_response:true }...: avg=162.64ms min=570.64µs med=134.47ms max=5.58s  p(90)=348.13ms p(95)=408.13ms
    http_req_failed................: 0.00%  30 out of 431748
    http_reqs......................: 431748 3568.709106/s

    EXECUTION
    iteration_duration.............: avg=758.48ms min=103.45ms med=655.63ms max=15.33s p(90)=1.52s    p(95)=1.69s
    iterations.....................: 107937 892.177276/s
    vus............................: 775    min=104          max=1494
    vus_max........................: 1500   min=1500         max=1500

    NETWORK
    data_received..................: 1.0 GB 8.4 MB/s
    data_sent......................: 113 MB 933 kB/s




running (2m01.0s), 0000/1500 VUs, 107937 complete and 0 interrupted iterations
stress_test ✓ [======================================] 0000/1500 VUs  2m0s
```

## User Spike Test

```text
 TOTAL RESULTS

    checks_total.......: 194508  3235.764502/s
    checks_succeeded...: 100.00% 194508 out of 194508
    checks_failed......: 0.00%   0 out of 194508

    ✓ GET /api/user?page=1&limit=10 success
    ✓ GET /api/user/41 success
    ✓ GET /api/user/active?page=1&limit=10 success
    ✓ GET /api/user/trashed?page=1&limit=10 success

    HTTP
    http_req_duration..............: avg=184.88ms min=567.94µs med=202ms    max=623.48ms p(90)=297.47ms p(95)=336.16ms
      { expected_response:true }...: avg=184.88ms min=567.94µs med=202ms    max=623.48ms p(90)=297.47ms p(95)=336.16ms
    http_req_failed................: 0.00%  0 out of 194508
    http_reqs......................: 194508 3235.764502/s

    EXECUTION
    iteration_duration.............: avg=843.33ms min=103.64ms med=963.35ms max=1.66s    p(90)=1.23s    p(95)=1.29s
    iterations.....................: 48627  808.941125/s
    vus............................: 82     min=4           max=1000
    vus_max........................: 1000   min=1000        max=1000

    NETWORK
    data_received..................: 455 MB 7.6 MB/s
    data_sent......................: 51 MB  846 kB/s




running (1m00.1s), 0000/1000 VUs, 48627 complete and 0 interrupted iterations
spike_test ✓ [======================================] 0000/1000 VUs  1m0s
```

---

## Role Capability

```text
█ TOTAL RESULTS

    checks_total.......: 216790  1203.6794/s
    checks_succeeded...: 100.00% 216790 out of 216790
    checks_failed......: 0.00%   0 out of 216790

    ✓ GET /api/role?page=1&limit=10 success
    ✓ GET /api/role/active?page=1&limit=10 success
    ✓ GET /api/role/trashed?page=1&limit=10 success
    ✓ GET /api/role/1 success
    ✓ GET /api/role/user/1 success

    HTTP
    http_req_duration..............: avg=2.14ms   min=505.37µs med=884.28µs max=116.39ms p(90)=3.38ms   p(95)=7.53ms
      { expected_response:true }...: avg=2.14ms   min=505.37µs med=884.28µs max=116.39ms p(90)=3.38ms   p(95)=7.53ms
    http_req_failed................: 0.00%  0 out of 216790
    http_reqs......................: 216790 1203.6794/s

    EXECUTION
    dropped_iterations.............: 138    0.766215/s
    iteration_duration.............: avg=112.24ms min=103.56ms med=105.88ms max=386.7ms  p(90)=119.04ms p(95)=143.03ms
    iterations.....................: 43358  240.73588/s
    vus............................: 63     min=5           max=128
    vus_max........................: 154    min=100         max=154

    NETWORK
    data_received..................: 191 MB 1.1 MB/s
    data_sent......................: 56 MB  312 kB/s




running (3m00.1s), 000/154 VUs, 43358 complete and 0 interrupted iterations
scalability_test ✓ [======================================] 000/154 VUs  3m0s  599.97 iters/s
```

## Role Load Test

```text
  http_req_duration
    ✗ 'p(95)<200' p(95)=333.09ms

    http_req_failed
    ✓ 'rate<0.01' rate=0.00%


  █ TOTAL RESULTS

    checks_total.......: 506945 4208.152943/s
    checks_succeeded...: 99.99% 506915 out of 506945
    checks_failed......: 0.00%  30 out of 506945

    ✗ GET /api/role?page=1&limit=10 success
      ↳  99% — ✓ 101361 / ✗ 28
    ✗ GET /api/role/active?page=1&limit=10 success
      ↳  99% — ✓ 101387 / ✗ 2
    ✓ GET /api/role/trashed?page=1&limit=10 success
    ✓ GET /api/role/1 success
    ✓ GET /api/role/user/1 success

    HTTP
    http_req_duration..............: avg=215.78ms min=5.24ms   med=197.59ms max=15.3s  p(90)=291.31ms p(95)=333.09ms
      { expected_response:true }...: avg=214.9ms  min=5.24ms   med=197.59ms max=5.74s  p(90)=291.29ms p(95)=333.03ms
    http_req_failed................: 0.00%  30 out of 506945
    http_reqs......................: 506945 4208.152943/s

    EXECUTION
    iteration_duration.............: avg=1.18s    min=457.26ms med=1.13s    max=16.55s p(90)=1.35s    p(95)=1.59s
    iterations.....................: 101389 841.630589/s
    vus............................: 1000   min=1000         max=1000
    vus_max........................: 1000   min=1000         max=1000

    NETWORK
    data_received..................: 446 MB 3.7 MB/s
    data_sent......................: 132 MB 1.1 MB/s




running (2m00.5s), 0000/1000 VUs, 101389 complete and 0 interrupted iterations
load_test ✓ [======================================] 1000 VUs  2m0s
ERRO[0121] thresholds on metrics 'http_req_duration' have been crossed
```

## Role Ramping Test

```text
 TOTAL RESULTS

    checks_total.......: 496950 3981.861522/s
    checks_succeeded...: 99.99% 496920 out of 496950
    checks_failed......: 0.00%  30 out of 496950

    ✓ GET /api/role?page=1&limit=10 success
    ✗ GET /api/role/active?page=1&limit=10 success
      ↳  99% — ✓ 99360 / ✗ 30
    ✓ GET /api/role/trashed?page=1&limit=10 success
    ✓ GET /api/role/1 success
    ✓ GET /api/role/user/1 success

    HTTP
    http_req_duration..............: avg=152.06ms min=554.06µs med=109.55ms max=15.1s  p(90)=269.49ms p(95)=349.33ms
      { expected_response:true }...: avg=151.16ms min=554.06µs med=109.54ms max=8.33s  p(90)=269.4ms  p(95)=349.17ms
    http_req_failed................: 0.00%  30 out of 496950
    http_reqs......................: 496950 3981.861522/s

    EXECUTION
    iteration_duration.............: avg=864.81ms min=104.8ms  med=667ms    max=15.39s p(90)=1.28s    p(95)=1.81s
    iterations.....................: 99390  796.372304/s
    vus............................: 195    min=104          max=1493
    vus_max........................: 1500   min=1500         max=1500

    NETWORK
    data_received..................: 437 MB 3.5 MB/s
    data_sent......................: 129 MB 1.0 MB/s




running (2m04.8s), 0000/1500 VUs, 99390 complete and 0 interrupted iterations
stress_test ✓ [======================================] 0000/1500 VUs  2m0s
```

### Role Spike Test

```text
 TOTAL RESULTS

    checks_total.......: 231610  3852.548945/s
    checks_succeeded...: 100.00% 231610 out of 231610
    checks_failed......: 0.00%   0 out of 231610

    ✓ GET /api/role?page=1&limit=10 success
    ✓ GET /api/role/active?page=1&limit=10 success
    ✓ GET /api/role/trashed?page=1&limit=10 success
    ✓ GET /api/role/1 success
    ✓ GET /api/role/user/1 success

    HTTP
    http_req_duration..............: avg=156.23ms min=545.93µs med=169.01ms max=554.66ms p(90)=245.17ms p(95)=270.87ms
      { expected_response:true }...: avg=156.23ms min=545.93µs med=169.01ms max=554.66ms p(90)=245.17ms p(95)=270.87ms
    http_req_failed................: 0.00%  0 out of 231610
    http_reqs......................: 231610 3852.548945/s

    EXECUTION
    iteration_duration.............: avg=886.06ms min=104.06ms med=1.02s    max=1.53s    p(90)=1.19s    p(95)=1.24s
    iterations.....................: 46322  770.509789/s
    vus............................: 84     min=4           max=1000
    vus_max........................: 1000   min=1000        max=1000

    NETWORK
    data_received..................: 204 MB 3.4 MB/s
    data_sent......................: 60 MB  999 kB/s




running (1m00.1s), 0000/1000 VUs, 46322 complete and 0 interrupted iterations
spike_test ✓ [======================================] 0000/1000 VUs  1m0s
hoover@hoover-IdeaPad-3-14IML05 ~/P/g/payment-gateway-grpc (main)>
```

---

## Card Capability

```text

  █ TOTAL RESULTS

    checks_total.......: 431040 2242.138799/s
    checks_succeeded...: 87.49% 377157 out of 431040
    checks_failed......: 12.50% 53883 out of 431040

    ✓ GET /api/card?page=1&limit=10 success
    ✓ GET /api/card/active?page=1&limit=10 success
    ✓ GET /api/card/trashed?page=1&limit=10 success
    ✗ GET /api/card/user?user_id=11 success
      ↳  0% — ✓ 0 / ✗ 13470
    ✓ GET /api/card/card_number/4389195775213720 success
    ✓ GET /api/card/11 success
    ✓ GET /api/card/dashboard success
    ✗ GET /api/card/dashboard/4016531504435114 success
      ↳  0% — ✓ 0 / ✗ 13470
    ✗ GET /api/card/monthly-balance?year=2025&month=1 success
      ↳  0% — ✓ 0 / ✗ 13470
    ✗ GET /api/card/yearly-balance?year=2025 success
      ↳  0% — ✓ 0 / ✗ 13470
    ✓ GET /api/card/monthly-balance-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-balance-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-topup-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-topup-amount?year=2025 success
    ✓ GET /api/card/monthly-topup-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-topup-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-withdraw-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-withdraw-amount?year=2025 success
    ✗ GET /api/card/monthly-withdraw-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
      ↳  99% — ✓ 13467 / ✗ 3
    ✓ GET /api/card/yearly-withdraw-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transaction-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transaction-amount?year=2025 success
    ✓ GET /api/card/monthly-transaction-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transaction-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-sender-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-sender-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-sender-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-sender-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-receiver-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-receiver-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-receiver-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-receiver-amount-by-card?year=2025&card_number=4016531504435114 success

    HTTP
    http_req_duration..............: avg=213.78ms min=608.07µs med=89.21ms max=17.37s p(90)=515.57ms p(95)=796.61ms
      { expected_response:true }...: avg=203.59ms min=608.07µs med=86.54ms max=5.11s  p(90)=507.75ms p(95)=783.39ms
    http_req_failed................: 12.50% 53883 out of 431040
    http_reqs......................: 431040 2242.138799/s

    EXECUTION
    dropped_iterations.............: 30020  156.154897/s
    iteration_duration.............: avg=6.95s    min=135.78ms med=4.08s   max=45.59s p(90)=18.87s   p(95)=25.5s
    iterations.....................: 13470  70.066837/s
    vus............................: 30     min=9               max=900
    vus_max........................: 900    min=100             max=900

    NETWORK
    data_received..................: 472 MB 2.5 MB/s
    data_sent......................: 126 MB 656 kB/s




running (3m12.2s), 000/900 VUs, 13470 complete and 0 interrupted iterations
scalability_test ✓ [======================================] 000/900 VUs  3m0s  599.92 iters/s
```

## Card Load Test

```text
 █ THRESHOLDS

    http_req_duration
    ✗ 'p(95)<200' p(95)=805.78ms

    http_req_failed
    ✗ 'rate<0.01' rate=12.51%


  █ TOTAL RESULTS

    checks_total.......: 374368 2601.312073/s
    checks_succeeded...: 87.48% 327517 out of 374368
    checks_failed......: 12.51% 46851 out of 374368

    ✗ GET /api/card?page=1&limit=10 success
      ↳  99% — ✓ 11669 / ✗ 30
    ✓ GET /api/card/active?page=1&limit=10 success
    ✓ GET /api/card/trashed?page=1&limit=10 success
    ✗ GET /api/card/user?user_id=11 success
      ↳  0% — ✓ 0 / ✗ 11699
    ✗ GET /api/card/card_number/4389195775213720 success
      ↳  99% — ✓ 11674 / ✗ 25
    ✓ GET /api/card/11 success
    ✓ GET /api/card/dashboard success
    ✗ GET /api/card/dashboard/4016531504435114 success
      ↳  0% — ✓ 0 / ✗ 11699
    ✗ GET /api/card/monthly-balance?year=2025&month=1 success
      ↳  0% — ✓ 0 / ✗ 11699
    ✗ GET /api/card/yearly-balance?year=2025 success
      ↳  0% — ✓ 0 / ✗ 11699
    ✓ GET /api/card/monthly-balance-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-balance-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-topup-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-topup-amount?year=2025 success
    ✓ GET /api/card/monthly-topup-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-topup-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-withdraw-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-withdraw-amount?year=2025 success
    ✓ GET /api/card/monthly-withdraw-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-withdraw-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transaction-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transaction-amount?year=2025 success
    ✓ GET /api/card/monthly-transaction-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transaction-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-sender-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-sender-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-sender-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-sender-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-receiver-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-receiver-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-receiver-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-receiver-amount-by-card?year=2025&card_number=4016531504435114 success

    HTTP
    http_req_duration..............: avg=352.93ms min=584.38µs med=270.43ms max=19.36s p(90)=576.6ms  p(95)=805.78ms
      { expected_response:true }...: avg=333.29ms min=584.38µs med=259.53ms max=6.6s   p(90)=554.29ms p(95)=776.46ms
    http_req_failed................: 12.51% 46851 out of 374368
    http_reqs......................: 374368 2601.312073/s

    EXECUTION
    iteration_duration.............: avg=11.41s   min=6.63s    med=8.76s    max=47.25s p(90)=20.99s   p(95)=23.45s
    iterations.....................: 11699  81.291002/s
    vus............................: 1      min=1               max=1000
    vus_max........................: 1000   min=1000            max=1000

    NETWORK
    data_received..................: 410 MB 2.8 MB/s
    data_sent......................: 110 MB 761 kB/s




running (2m23.9s), 0000/1000 VUs, 11699 complete and 0 interrupted iterations
load_test ✓ [======================================] 1000 VUs  2m0s
ERRO[0144] thresholds on metrics 'http_req_duration, http_req_failed' have been crossed
```

## Card Rampling Test

```text

  █ TOTAL RESULTS

    checks_total.......: 415680 3232.505582/s
    checks_succeeded...: 87.49% 363690 out of 415680
    checks_failed......: 12.50% 51990 out of 415680

    ✗ GET /api/card?page=1&limit=10 success
      ↳  99% — ✓ 12960 / ✗ 30
    ✓ GET /api/card/active?page=1&limit=10 success
    ✓ GET /api/card/trashed?page=1&limit=10 success
    ✗ GET /api/card/user?user_id=11 success
      ↳  0% — ✓ 0 / ✗ 12990
    ✓ GET /api/card/card_number/4389195775213720 success
    ✓ GET /api/card/11 success
    ✓ GET /api/card/dashboard success
    ✗ GET /api/card/dashboard/4016531504435114 success
      ↳  0% — ✓ 0 / ✗ 12990
    ✗ GET /api/card/monthly-balance?year=2025&month=1 success
      ↳  0% — ✓ 0 / ✗ 12990
    ✗ GET /api/card/yearly-balance?year=2025 success
      ↳  0% — ✓ 0 / ✗ 12990
    ✓ GET /api/card/monthly-balance-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-balance-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-topup-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-topup-amount?year=2025 success
    ✓ GET /api/card/monthly-topup-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-topup-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-withdraw-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-withdraw-amount?year=2025 success
    ✓ GET /api/card/monthly-withdraw-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-withdraw-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transaction-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transaction-amount?year=2025 success
    ✓ GET /api/card/monthly-transaction-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transaction-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-sender-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-sender-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-sender-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-sender-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-receiver-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-receiver-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-receiver-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-receiver-amount-by-card?year=2025&card_number=4016531504435114 success

    HTTP
    http_req_duration..............: avg=210.98ms min=644.91µs med=143.99ms max=15.64s p(90)=476.95ms p(95)=586.85ms
      { expected_response:true }...: avg=201.17ms min=644.91µs med=141.58ms max=5.3s   p(90)=469.56ms p(95)=578.87ms
    http_req_failed................: 12.50% 51990 out of 415680
    http_reqs......................: 415680 3232.505582/s

    EXECUTION
    iteration_duration.............: avg=6.86s    min=379.67ms med=4.86s    max=42.56s p(90)=14.89s   p(95)=16.48s
    iterations.....................: 12990  101.015799/s
    vus............................: 29     min=29              max=1494
    vus_max........................: 1500   min=1500            max=1500

    NETWORK
    data_received..................: 455 MB 3.5 MB/s
    data_sent......................: 122 MB 946 kB/s




running (2m08.6s), 0000/1500 VUs, 12990 complete and 0 interrupted iterations
stress_test ✓ [======================================] 0000/1500 VUs  2m0s
```

## Card Spike Test

```text
 █ TOTAL RESULTS

    checks_total.......: 201760 3088.166155/s
    checks_succeeded...: 87.50% 176540 out of 201760
    checks_failed......: 12.50% 25220 out of 201760

    ✓ GET /api/card?page=1&limit=10 success
    ✓ GET /api/card/active?page=1&limit=10 success
    ✓ GET /api/card/trashed?page=1&limit=10 success
    ✗ GET /api/card/user?user_id=11 success
      ↳  0% — ✓ 0 / ✗ 6305
    ✓ GET /api/card/card_number/4389195775213720 success
    ✓ GET /api/card/11 success
    ✓ GET /api/card/dashboard success
    ✗ GET /api/card/dashboard/4016531504435114 success
      ↳  0% — ✓ 0 / ✗ 6305
    ✗ GET /api/card/monthly-balance?year=2025&month=1 success
      ↳  0% — ✓ 0 / ✗ 6305
    ✗ GET /api/card/yearly-balance?year=2025 success
      ↳  0% — ✓ 0 / ✗ 6305
    ✓ GET /api/card/monthly-balance-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-balance-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-topup-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-topup-amount?year=2025 success
    ✓ GET /api/card/monthly-topup-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-topup-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-withdraw-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-withdraw-amount?year=2025 success
    ✓ GET /api/card/monthly-withdraw-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-withdraw-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transaction-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transaction-amount?year=2025 success
    ✓ GET /api/card/monthly-transaction-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transaction-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-sender-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-sender-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-sender-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-sender-amount-by-card?year=2025&card_number=4016531504435114 success
    ✓ GET /api/card/monthly-transfer-receiver-amount?year=2025&month=1 success
    ✓ GET /api/card/yearly-transfer-receiver-amount?year=2025 success
    ✓ GET /api/card/monthly-transfer-receiver-amount-by-card?year=2025&month=1&card_number=4016531504435114 success
    ✓ GET /api/card/yearly-transfer-receiver-amount-by-card?year=2025&card_number=4016531504435114 success

    HTTP
    http_req_duration..............: avg=212.85ms min=580.78µs med=241.75ms max=15.3s  p(90)=330.13ms p(95)=359.86ms
      { expected_response:true }...: avg=207.95ms min=580.78µs med=239.7ms  max=1.16s  p(90)=327.55ms p(95)=356.83ms
    http_req_failed................: 12.50% 25220 out of 201760
    http_reqs......................: 201760 3088.166155/s

    EXECUTION
    iteration_duration.............: avg=6.92s    min=133.93ms med=8.53s    max=17.65s p(90)=8.98s    p(95)=9.11s
    iterations.....................: 6305   96.505192/s
    vus............................: 30     min=4               max=1000
    vus_max........................: 1000   min=1000            max=1000

    NETWORK
    data_received..................: 221 MB 3.4 MB/s
    data_sent......................: 59 MB  904 kB/s




running (1m05.3s), 0000/1000 VUs, 6305 complete and 0 interrupted iterations
spike_test ✓ [======================================] 0000/1000 VUs  1m0s
```

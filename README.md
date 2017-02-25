![mercury](https://github.com/ortuman/mercury/blob/master/screenshots/logo.png)

## Installation

```bash
go get -u github.com/ortuman/mercury
```

## Support Platform

* [APNS](https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/APNSOverview.html#//apple_ref/doc/uid/TP40008194-CH8-SW1)
* [GCM](https://developer.android.com/google/gcm/index.html)
* [Safari](https://developer.apple.com/library/prerelease/content/documentation/NetworkingInternet/Conceptual/NotificationProgrammingGuideForWebsites/PushNotifications/PushNotifications.html)
* [Chrome](https://developers.google.com/web/fundamentals/getting-started/codelabs/push-notifications/)
* [Firefox](https://support.mozilla.org/t5/Basic-Browsing/Web-Push-notifications-in-Firefox/ta-p/28744)

## Web API

### GET /ping

Checks if the server is alive. This is useful for monitoring tools, load-balancers and automated scripts.

### GET /v1/stats/push

Show push notification stats.

```json
{
  "apns": {
    "max_conn": 16,
    "delivered_count": 243386,
    "unregistered_count": 18608,
    "failed_count": 8,
    "avg_request_time_msec": 8
  },
  "chrome": {
    "max_conn": 16,
    "delivered_count": 196,
    "unregistered_count": 18,
    "failed_count": 0,
    "avg_request_time_msec": 32
  },
  "firefox": {
    "max_conn": 16,
    "delivered_count": 323,
    "unregistered_count": 21,
    "failed_count": 0,
    "avg_request_time_msec": 17
  },
  "gcm": {
    "max_conn": 16,
    "delivered_count": 0,
    "unregistered_count": 0,
    "failed_count": 3,
    "avg_request_time_msec": 12
  },
  "safari": {
    "max_conn": 8,
    "delivered_count": 504,
    "unregistered_count": 7,
    "failed_count": 12,
    "avg_request_time_msec": 2
  }
}
```

### GET /v1/stats/sys

Show system stats.

```json
{
    "time": 1488047755861209300,
    "go_version": "go1.8",
    "go_os": "darwin",
    "go_arch": "amd64",
    "cpu_num": 4,
    "goroutine_num": 4,
    "gomaxprocs": 4,
    "cgo_call_num": 1,
    "memory_alloc": 1797408,
    "memory_total_alloc": 1797408,
    "memory_sys": 5343480,
    "memory_lookups": 8,
    "memory_mallocs": 12897,
    "memory_frees": 407,
    "memory_stack": 327680,
    "heap_alloc": 1797408,
    "heap_sys": 2818048,
    "heap_idle": 466944,
    "heap_inuse": 2351104,
    "heap_released": 0,
    "heap_objects": 12490,
    "gc_next": 4194304,
    "gc_last": 1460686815762559700,
    "gc_num": 1,
    "gc_per_second": 0,
    "gc_pause_per_second": 0,
    "gc_pause": [ 0.395676 ]
}
```

## License

Copyright 2017 Miguel Ángel Ortuño <ortuman@gmail.com>

Licensed under the MIT License.

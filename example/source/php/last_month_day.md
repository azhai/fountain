---
date:       2016-03-24
title:      找出上个月的这一天
slug:       last-month-day
author:     ryan
tags:       [php, date, month]
comments:   true
weight:     560324
---

找出上个月的这一天，没有这一天时使用月末

```php
/**
 * 找出上个月的这一天，没有这一天时使用月末
 * （使用时间戳计算，避免判断跨年）
 */
function last_month_day($time)
{
    $day = intval(date('d', $time)); //当月第几天
    $time -= $day * 86400; //退回上月最后一天
    $tail_day = intval(date('d', $time)); //上个月有多长
    if ($day > $tail_day) {
        //上个月较短，没有这几天，使用月末
    } else {
        $time -= ($tail_day - $day) * 86400;
    }
    return date('Y-m-d', $time);
}

$time = strtotime('2012-03-31');
echo date('Y-m-d', $time) . "\n";
echo last_month_day($time) . "\n";
```

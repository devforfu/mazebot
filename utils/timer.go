package utils

import "time"

func Timer(function func()()) time.Duration {
    t := time.Now()
    function()
    elapsed := time.Since(t)
    return elapsed
}

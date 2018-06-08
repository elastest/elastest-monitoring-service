package datatypes

func Min(t0 MaybeTime, t1 MaybeTime) MaybeTime {
    if !t0.IsSet {
        return t1
    }
    if !t1.IsSet {
        return t0
    }
    if t0.Val < t1.Val {
        return t0
    }
    return t1
}

func Leq(t0 Time, t1 Time) bool {
    return t0 <= t1
}

func Lt(t0 Time, t1 Time) bool {
    return t0 < t1
}

func FstPayload(t0 EvPayload, _ EvPayload) EvPayload {
    return t0
}

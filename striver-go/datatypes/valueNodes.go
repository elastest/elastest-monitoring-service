package datatypes

// TNode

func (node TNode) Exec (t Time, _ interface{}, _ InPipes) EvPayload {
    return Some(t)
}

func (node TNode) Rinse (_ InPipes) {
}

// WNode

func (node WNode) Exec (_ Time, w interface{}, _ InPipes) EvPayload {
    return Some(w)
}

func (node WNode) Rinse (inpipes InPipes) {
}

// aux funs

func consumeWhile(seen []Event, cmpfun func(Time) bool) ([]Event, *Event) {
    if (len(seen) == 0 || !cmpfun(seen[0].Time)) {
        return seen, nil // outside
    }
    i:=1
    for ;i<len(seen);i++ {
        if !cmpfun(seen[i].Time) {
            break
        }
    }
    seen = seen[i-1:]
    return seen, &seen[0]
}

func norinse(InPipes) {}

// Beta testing: generic funs

func genericExec(t Time, w interface{}, inpipes InPipes, rinsefun func(InPipes), tpointernode ValNode, cmpfun func(t0 Time, t1 Time) bool, seen *[]Event, extractor func(Event) interface{}) EvPayload {
    tpayload := tpointernode.Exec(t,w,inpipes)
    if !tpayload.IsSet {
        // outside
        return tpayload
    }
    limitT := tpayload.Val.(Time)
    if (limitT == t) {
        // It might be now
        rinsefun(inpipes)
    }
    newseen, ev := consumeWhile(*seen, func(t Time) bool {return cmpfun(t, limitT)})
    *seen = newseen
    if ev==nil {
        // outside
        return NothingPayload
    }
    if (!ev.Payload.IsSet) {
        panic("Empty payload in queue??")
    }
    return Some(extractor(*ev))
}

func genericRinse (inpipes InPipes, tpointernode ValNode, srcStream StreamName, seen *[]Event) {
    tpointernode.Rinse(inpipes)
    ev := inpipes.strictConsume(srcStream)
    if ev.Payload.IsSet && (len(*seen) == 0 || ev.Time!=(*seen)[len(*seen)-1].Time) {
        *seen = (append(*seen, ev))
    }
}

func extractPayload(ev Event) interface{} {
    return ev.Payload
}

func extractTime(ev Event) interface{} {
    return ev.Time
}

// PrevEqValNode

func (node *PrevEqValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return genericExec(t, w, inpipes, node.Rinse, node.TPointer, Leq, &node.Seen, extractPayload)
}

func (node *PrevEqValNode) Rinse (inpipes InPipes) {
    genericRinse(inpipes, node.TPointer, node.SrcStream, &node.Seen)
}

// PrevValNode

func (node *PrevValNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return genericExec(t, w, inpipes, norinse, node.TPointer, Lt, &node.Seen, extractPayload)
}

func (node *PrevValNode) Rinse (inpipes InPipes) {
    genericRinse(inpipes, node.TPointer, node.SrcStream, &node.Seen)
}

// PrevEqNode

func (node *PrevEqNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return genericExec(t, w, inpipes, node.Rinse, node.TPointer, Leq, &node.Seen, extractTime)
}

func (node *PrevEqNode) Rinse (inpipes InPipes) {
    genericRinse(inpipes, node.TPointer, node.SrcStream, &node.Seen)
}

// PrevNode

func (node *PrevNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    return genericExec(t, w, inpipes, norinse, node.TPointer, Lt, &node.Seen, extractTime)
}

func (node *PrevNode) Rinse (inpipes InPipes) {
    genericRinse(inpipes, node.TPointer, node.SrcStream, &node.Seen)
}

// FuncNode

func (node FuncNode) Exec (t Time, w interface{}, inpipes InPipes) EvPayload {
    args := make([]EvPayload, len(node.ArgNodes))
    for i,valnode := range node.ArgNodes {
        args[i] = valnode.Exec(t, w, inpipes)
    }
    return node.Innerfun(args...)
}

func (node FuncNode) Rinse (inpipes InPipes) {
    for _,valnode := range node.ArgNodes {
        valnode.Rinse(inpipes)
    }
}

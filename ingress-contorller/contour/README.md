# contour:envoy-ingress-controller
### github repo:https://github.com/heptio/contour

### 严重风险-1
在用户配置 spec.rules.host 为空时,会配置默认 host = *
所有流量会被引入这个可以接受任何host的rule，造成服务异常
contour _contour/internal/contour/dag/dag.go:320_
```
for _, rule := range ing.Spec.Rules {
    // handle Spec.Rule declarations
    host := rule.Host
    if host == "" {
        host = "*"
    }
```

### 严重风险-2
在用户配置 spec.rule.http及其它选项都为空时，&Route会指向一个无效地址，导致contour panic，影响整个控制层
contour _contour/internal/contour/dag/dag.go:331_
```
				r := &Route{
					path:         path,
					object:       ing,
					HTTPSUpgrade: tlsRequired(ing),
					Websocket:    wr[path],
					Timeout:      timeout,
				}
```
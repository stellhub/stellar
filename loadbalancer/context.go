package loadbalancer

import "context"

type requestContextKey struct{}

func ContextWithRequest(ctx context.Context, request Request) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, requestContextKey{}, request)
}

func RequestFromContext(ctx context.Context) (Request, bool) {
	if ctx == nil {
		return Request{}, false
	}
	request, ok := ctx.Value(requestContextKey{}).(Request)
	return request, ok
}

func MergeContextRequest(ctx context.Context, request Request) Request {
	if existing, ok := RequestFromContext(ctx); ok {
		request = mergeRequest(existing, request)
	}
	return request
}

func mergeRequest(base Request, override Request) Request {
	if override.Protocol != "" {
		base.Protocol = override.Protocol
	}
	if override.Service != "" {
		base.Service = override.Service
	}
	if override.Operation != "" {
		base.Operation = override.Operation
	}
	if override.Method != "" {
		base.Method = override.Method
	}
	if override.Path != "" {
		base.Path = override.Path
	}
	if override.Target != "" {
		base.Target = override.Target
	}
	if override.HashKey != "" {
		base.HashKey = override.HashKey
	}
	if override.Headers != nil {
		base.Headers = override.Headers
	}
	if override.Attributes != nil {
		if base.Attributes == nil {
			base.Attributes = map[string]any{}
		}
		for key, value := range override.Attributes {
			base.Attributes[key] = value
		}
	}
	return base
}

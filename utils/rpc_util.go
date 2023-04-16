package utils

import (
	"encoding/json"
	"gogo/constant"
	"io"
	"log"
	"net/http"
	"strings"
)

var HeadersToPropagate = []string{
	// All applications should propagate x-request-id. This header is
	// included in access log statements and is used for consistent trace
	// sampling and log sampling decisions in Istio.
	"x-request-id",

	// Lightstep tracing header. Propagate this if you use lightstep tracing
	// in Istio (see
	// https://istio.io/latest/docs/tasks/observability/distributed-tracing/lightstep/)
	// Note: this should probably be changed to use B3 or W3C TRACE_CONTEXT.
	// Lightstep recommends using B3 or TRACE_CONTEXT and most application
	// libraries from lightstep do not support x-ot-span-context.
	"x-ot-span-context",

	// Datadog tracing header. Propagate these headers if you use Datadog
	// tracing.
	"x-datadog-trace-id",
	"x-datadog-parent-id",
	"x-datadog-sampling-priority",

	// W3C Trace Context. Compatible with OpenCensusAgent and Stackdriver Istio
	// configurations.
	"traceparent",
	"tracestate",

	// Cloud trace context. Compatible with OpenCensusAgent and Stackdriver Istio
	// configurations.
	"x-cloud-trace-context",

	// Grpc binary trace context. Compatible with OpenCensusAgent nad
	// Stackdriver Istio configurations.
	"grpc-trace-bin",

	// b3 trace headers. Compatible with Zipkin, OpenCensusAgent, and
	// Stackdriver Istio configurations. Commented out since they are
	// propagated by the OpenTracing tracer above.
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",

	// SkyWalking trace headers.
	"sw8",

	// Application-specific headers to forward.
	"end-user",
	"user-agent",

	// Context and session specific headers
	"cookie",
	"authorization",
	"jwt",
}

type PageRsp struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Pages int64       `json:"pages"`
}

func RpcGet(appName string, path string, params string, header http.Header) *PageRsp {
	return RpcDoWithBody(appName, path, params, header, http.MethodGet)
}

func RpcPost(appName string, path string, params string, header http.Header) *PageRsp {
	return RpcDoWithBody(appName, path, params, header, http.MethodPost)
}

func RpcPut(appName string, path string, params string, header http.Header) *PageRsp {
	return RpcDoWithBody(appName, path, params, header, http.MethodPut)
}

func RpcDelete(appName string, path string, params string, header http.Header) *PageRsp {
	return RpcDoWithBody(appName, path, params, header, http.MethodDelete)
}

func RpcDoWithBody(appName string, path string, params string, header http.Header, requestMethod string) *PageRsp {
	//http://resource.gogo-system.svc.cluster.local:9080/resource/health
	url := "http://" + appName + constant.Dot + constant.Namespace + constant.Dot + "svc.cluster.local" + constant.ReleaseAddr + path
	//url := "http://" + appName + "/dir/dirRoot/getRoot"
	req, err := http.NewRequest(http.MethodDelete, url, strings.NewReader(params))
	for _, key := range HeadersToPropagate {
		if IsNotEmpty(header[key]) && IsNotBlack(header[key][0]) {
			req.Header.Set(key, header[key][0])
		}
	}
	//req.Header.Set("ContentType","application/gprc")
	if err != nil {
		log.Printf("RpcPost error: %v", err)
	}
	//tlsConfig := &tls.Config{}
	//transport := &http2.Transport{TLSClientConfig: tlsConfig}
	client := http.Client{}
	println(url)
	rsp, err := client.Do(req)
	if err != nil {
		log.Panicf("RpcPost error: %v", err)
	}
	defer rsp.Body.Close()
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("RpcPost error: %v", err)
	}
	r := &PageRsp{}
	err = json.Unmarshal(body, r)
	if err != nil {
		log.Panicln("返回参数不是有效的Rsp格式")
	}
	return r
}

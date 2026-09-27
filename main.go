package main

import (
    "fmt"
    "net/http"
    "os"
    "math"
    "strconv"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var version = "dev"

var workRequests = prometheus.NewCounter(prometheus.CounterOpts{
    Name: "testapp_work_requests_total",
    Help: "Number of /work requests.",
})

func healthz(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "go-app version: %s\n", version)
}

func work(w http.ResponseWriter, r *http.Request) {
    workRequests.Inc()
    n := 1000000
    if q := r.URL.Query().Get("n"); q != "" {
        if v, err := strconv.Atoi(q); err == nil && v > 0 {
            n = v
        }
    }
    x := 0.5
    for i := 0; i < n; i++ {
        x = math.Sqrt(x*x + 2)
    }
    fmt.Fprintf(w, "work: %d iterations = %f\n", n, x)
}

func init() {
    prometheus.MustRegister(workRequests)
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/", handler)
    http.HandleFunc("/healthz", healthz)
    http.HandleFunc("/work", work)

    metricsMux := http.NewServeMux()
    metricsMux.Handle("/metrics", promhttp.Handler())
    go func() {
        if err := http.ListenAndServe(":9090", metricsMux); err != nil {
            fmt.Fprintln(os.Stderr, "metrics:", err)
        }
    }()

    fmt.Printf("listening on :%s (version %s)\n", port, version)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

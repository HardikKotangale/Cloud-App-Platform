package cli

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"

	"github.com/HardikKotangale/Cloud-App-Platform/internal/observability"
)

func NewMetricsCmd() *cobra.Command {
	var port int

	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Run a local /metrics endpoint (Prometheus format) for platform observability",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create a registry (so we control what gets exposed)
			reg := prometheus.NewRegistry()

			// Expose persistent totals as Gauges (simple & beginner-friendly)
			deploySuccess := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "platform_deploy_success_total",
				Help: "Total successful deployments performed via platformctl",
			}, func() float64 {
				c, _ := observability.Load()
				return float64(c.DeploySuccess)
			})

			deployFailure := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "platform_deploy_failure_total",
				Help: "Total failed deployments performed via platformctl",
			}, func() float64 {
				c, _ := observability.Load()
				return float64(c.DeployFailure)
			})

			policyViolations := prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "platform_policy_violations_total",
				Help: "Total governance/policy violations detected by platformctl",
			}, func() float64 {
				c, _ := observability.Load()
				return float64(c.PolicyViolations)
			})

			reg.MustRegister(deploySuccess, deployFailure, policyViolations)

			mux := http.NewServeMux()
			mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
			mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok"))
			})

			addr := fmt.Sprintf(":%d", port)
			srv := &http.Server{
				Addr:              addr,
				Handler:           mux,
				ReadHeaderTimeout: 5 * time.Second,
			}

			fmt.Printf("📈 Metrics server running:\n")
			fmt.Printf("  - http://localhost:%d/metrics\n", port)
			fmt.Printf("  - http://localhost:%d/healthz\n", port)
			fmt.Println("\nTip: In another terminal run: curl -s http://localhost:" + fmt.Sprint(port) + "/metrics | head")

			return srv.ListenAndServe()
		},
	}

	cmd.Flags().IntVar(&port, "port", 9090, "Port to listen on")
	return cmd
}

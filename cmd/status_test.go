package cmd

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"
)

func TestStatus_ConnectionRefused(t *testing.T) {

	cmd := exec.Command("go", "run", "../main.go", "status", "--region", "foobar", "--url", "https://localhost")
	out, _ := cmd.CombinedOutput()

	actual := string(out)
	expected := "[UNKNOWN] - Get \"https://localhost"

	if !strings.Contains(actual, expected) {
		t.Error("\nActual: ", actual, "\nExpected: ", expected)
	}
}

type StatusTest struct {
	name     string
	server   *httptest.Server
	args     []string
	expected string
}

func TestStatusCmd(t *testing.T) {
	tests := []StatusTest{
		{
			name: "status-unknown",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Nothing so split into Slice]]></title>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[UNKNOWN] - Could not determine status. Nothing so split into Slice",
		},
		{
			name: "status-ok",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Informational message: [RESOLVED] Elevated request error rate using the PUT object]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	 <item>
		<title><![CDATA[Informational message: [RESOLVED] Foobar Event]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 20 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[OK] - Event resolved for ec2 (eu-foobar-1)",
		},
		{
			name: "status-ok-no-items",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[OK] - No events for ec2 (eu-foobar-1)",
		},
		{
			name: "status-ok-normal",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Service is operating normally: Nothing to see move along]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[OK] - Service ec2 is operating normally (eu-foobar-1)",
		},
		{
			name: "status-warning",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Performance issues: Slow news day]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[WARNING] - Performance issues for ec2 (eu-foobar-1)",
		},
		{
			name: "status-critical",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Service disruption: Oh no!]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[CRITICAL] - Service disruption for ec2 (eu-foobar-1)",
		},
		{
			name: "status-informational",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Informational message: Foobar Event is unresolved]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "--region", "eu-foobar-1"},
			expected: "[WARNING] - Information available for ec2 (eu-foobar-1)",
		},
		{
			name: "status-global-service",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title><![CDATA[Amazon Simple Storage Service (Foobar) Service Status]]></title>
		<link>http://httptest.localhost/</link>
		<language>en-us</language>
		<lastBuildDate>Mon, 02 Jan 2023 00:05:49 PST</lastBuildDate>
		<generator>AWS Service Health Dashboard RSS Generator</generator>
		<description><![CDATA[Amazon  EventBridge Scheduler (Foobar) Service Status]]></description>
		<ttl>5</ttl>
	 <item>
		<title><![CDATA[Informational message: Foobar Event is unresolved]]></title>
		<link>http://httptest.localhost/</link>
		<pubDate>Fri, 24 Jul 2015 11:54:38 PDT</pubDate>
		<guid isPermaLink="false">http://httptest.localhost/1234</guid>
	 </item>
	</channel>
</rss>`))
			})),
			args:     []string{"run", "../main.go", "status", "-s", "iam", "--region", ""},
			expected: "[WARNING] - Information available for iam (Global)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer test.server.Close()

			cmd := exec.Command("go", append(test.args, "--url", test.server.URL)...)
			out, _ := cmd.CombinedOutput()

			actual := string(out)

			if !strings.Contains(actual, test.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", test.expected)
			}

		})
	}
}

package hcloud

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hetznercloud/hcloud-go/v2/hcloud/schema"
)

func TestPricingClientGet(t *testing.T) {
	env := newTestEnv()
	defer env.Teardown()

	env.Mux.HandleFunc("/pricing", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(schema.PricingGetResponse{
			Pricing: schema.Pricing{
				Currency: "EUR",
			},
		})
	})
	ctx := context.Background()

	pricing, _, err := env.Client.Pricing.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if pricing.Image.PerGBMonth.Currency != "EUR" {
		t.Errorf("unexpected currency: %v", pricing.Image.PerGBMonth.Currency)
	}
}

func TestPricingSchema(t *testing.T) {
	data := []byte(`{
		"currency": "EUR",
		"vat_rate": "19.00",
		"image": {
			"price_per_gb_month": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"floating_ip": {
			"price_monthly": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"floating_ips": [
			  {
				"prices": [
				  {
					"location": "fsn1",
					"price_monthly": {
					  "gross": "1.19",
					  "net": "1"
					}
				  }
				],
				"type": "ipv4"
			  }
			],
		"primary_ips": [
			{
				"prices": [
				{
					"location": "fsn1",
					"price_hourly": {
					"gross": "1.1900000000000000",
					"net": "1.0000000000"
					},
					"price_monthly": {
					"gross": "1.1900000000000000",
					"net": "1.0000000000"
					}
				}
				],
				"type": "ipv4"
			}
		],
		"traffic": {
			"price_per_tb": {
				"net": "1",
				"gross": "1.19"
			}
		},
		"server_backup": {
			"percentage": "20"
		},
		"server_types": [
			{
				"id": 4,
				"name": "CX11",
				"prices": [
					{
						"location": "fsn1",
						"price_hourly": {
							"net": "1",
							"gross": "1.19"
						},
						"price_monthly": {
							"net": "1",
							"gross": "1.19"
						}
					}
				]
			}
		],
		"load_balancer_types": [
			{
				"id": 4,
				"name": "LX11",
				"prices": [
					{
						"location": "fsn1",
						"price_hourly": {
							"net": "1",
							"gross": "1.19"
						},
						"price_monthly": {
							"net": "1",
							"gross": "1.19"
						}
					}
				]
			}
		],
		"volume": {
			"price_per_gb_month": {
				"net": "1",
				"gross": "1.19"
			}
		}
	}`)

	var s schema.Pricing
	assert.NoError(t, json.Unmarshal(data, &s))

	assert.Equal(t, s, SchemaFromPricing(PricingFromSchema(s)))

	p := PricingFromSchema(s)
	assert.Equal(t, p, PricingFromSchema(SchemaFromPricing(p)))
}

func TestPricingSchemaPrimaryIPs(t *testing.T) {
	prices := `[{"location": "fsn1", "price_hourly": {"net": "1", "gross": "1.19"}, "price_monthly": {"net": "1", "gross": "1.19"}}]`

	t.Run("more primary IP types than floating IP types", func(t *testing.T) {
		data := []byte(`{
			"floating_ips": [{"type": "ipv4", "prices": ` + prices + `}],
			"primary_ips": [{"type": "ipv4", "prices": ` + prices + `}, {"type": "ipv6", "prices": ` + prices + `}]
		}`)
		var s schema.Pricing
		assert.NoError(t, json.Unmarshal(data, &s))

		p := PricingFromSchema(s)
		if assert.Len(t, p.PrimaryIPs, 2) {
			assert.Equal(t, "ipv4", p.PrimaryIPs[0].Type)
			assert.Equal(t, "ipv6", p.PrimaryIPs[1].Type)
		}
	})

	t.Run("no floating IP types", func(t *testing.T) {
		data := []byte(`{"primary_ips": [{"type": "ipv4", "prices": ` + prices + `}]}`)
		var s schema.Pricing
		assert.NoError(t, json.Unmarshal(data, &s))

		p := PricingFromSchema(s)
		if assert.Len(t, p.PrimaryIPs, 1) {
			assert.Equal(t, "1", p.PrimaryIPs[0].Pricings[0].Monthly.Net)
		}
	})

	t.Run("more floating IP types than primary IP types", func(t *testing.T) {
		data := []byte(`{
			"floating_ips": [{"type": "ipv4", "prices": ` + prices + `}, {"type": "ipv6", "prices": ` + prices + `}],
			"primary_ips": [{"type": "ipv4", "prices": ` + prices + `}]
		}`)
		var s schema.Pricing
		assert.NoError(t, json.Unmarshal(data, &s))

		assert.Len(t, PricingFromSchema(s).PrimaryIPs, 1)
	})
}

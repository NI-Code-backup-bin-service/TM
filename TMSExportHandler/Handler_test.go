package TMSExportHandler

import (
	_ "bitbucket.org/network-international/nextgen-libs/nextgen-helpers/exportHandler"
	"nextgen-tms-website/PED"
	"reflect"
	"testing"
)

var defaultMockData = []*PED.PEDDetailed{
	{
		PED: PED.PED{
			TID:      55001570,
			Serial:   "1502010369",
			SiteId:   88882973,
			SiteName: "PED Auto Test Site",
			PIN:      "1234",
		},
		Mode:        "standalone",
		Name:        "PED Auto Test Site",
		Active:      "[\"sale\",\"refund\",\"void\",\"preAuth\",\"gratuitySale\",\"gratuityCompletion\",\"alipay\",\"upi\",\"xls\",\"visaQr\",\"mastercardQr\",\"eppVoid\",\"balanceInquiry\"]",
		GratuityMax: 20,
	},
	{
		PED: PED.PED{
			TID:      88884574,
			Serial:   "1502000436",
			SiteId:   88882862,
			SiteName: "Lil Leigh's Bakery",
		},
		Mode:        "standalone",
		Name:        "Lil Leigh's Bakery",
		Active:      "[\"sale\",\"refund\",\"void\"]",
		GratuityMax: 0,
	},
}

type mockPEDRepository struct {
	pedDetails []*PED.PEDDetailed
}

func (r *mockPEDRepository) DeleteByTid(tid string) (tidDeleted bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *mockPEDRepository) DeleteOverrideByTid(tid string) (overrideDeleted bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *mockPEDRepository) DeleteFraudOverrideByTid(tid string) (fraudOverrideDeleted bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *mockPEDRepository) DeleteUserOverrideByTid(tid string) (userOverrideDeleted bool, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *mockPEDRepository) FindBySearchTermAndAcquirer(searchTerm string, acquirers string) ([]*PED.PEDDetailed, error) {
	return r.pedDetails, nil
}

func TestNewHandler(t *testing.T) {
	type args struct {
		repository PED.Repository
	}
	tests := []struct {
		name string
		args args
		want Handler
	}{
		{
			"Test1",
			args{
				&mockPEDRepository{
					defaultMockData,
				},
			},
			&handler{
				&mockPEDRepository{
					defaultMockData,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHandler(tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

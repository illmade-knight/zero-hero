package engine

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFromJSON(t *testing.T) {
	type args struct {
		mlp    string
		inputs [][]*Value
	}
	tests := []struct {
		name    string
		args    args
		want    [][]float64
		wantErr bool
	}{
		{"test load",
			args{
				"{\"inputs\":3,\"outputs\":1,\"layers\":[{\"neurons\":[{\"weights\":[{\"uuid\":\"39447886-203f-4cba-bb91-a00969f32e97\",\"label\":\"w-LMxV\",\"data\":-0.814235091944207},{\"uuid\":\"6ca531d5-1bc2-4fd7-ab65-3a224bbb4f81\",\"label\":\"w-esev\",\"data\":-0.8450130473625359},{\"uuid\":\"c57468f7-52d3-4218-8b39-87fe7e62b9c8\",\"label\":\"w-mZpL\",\"data\":-0.06818668325542561}],\"bias\":{\"uuid\":\"b037f12e-0293-4b55-b3bd-ddd78c705993\",\"label\":\"b-edmL\",\"data\":0.7218348447525101}},{\"weights\":[{\"uuid\":\"31fd8c8f-1574-43a0-a922-faea500570ae\",\"label\":\"w-nIhW\",\"data\":-0.45967633885311737},{\"uuid\":\"f88bac84-1fec-47c2-971d-6fe663941082\",\"label\":\"w-VJLH\",\"data\":-0.2740894934102625},{\"uuid\":\"ffb66bdb-71f9-457e-876d-10f4d5effd47\",\"label\":\"w-ENSi\",\"data\":-0.9347270788205421}],\"bias\":{\"uuid\":\"e4eb0755-8a12-4403-9dd2-74a647df3654\",\"label\":\"b-DLNu\",\"data\":-0.5933431363263343}},{\"weights\":[{\"uuid\":\"225d1afa-b330-4ab8-a8a8-0d41e3087fc9\",\"label\":\"w-rRiE\",\"data\":-0.9524115320261372},{\"uuid\":\"41311442-afcb-4162-af85-d12cc0dcc856\",\"label\":\"w-oKlw\",\"data\":0.585790427979633},{\"uuid\":\"68210d37-efce-45f0-a81c-7535f6e5c360\",\"label\":\"w-uYRc\",\"data\":-0.9145974669138838}],\"bias\":{\"uuid\":\"4211574e-42c8-473f-9000-2736006ca366\",\"label\":\"b-Ruzr\",\"data\":-0.023395312214124234}},{\"weights\":[{\"uuid\":\"6622224d-a876-4e36-b3e8-ce7b1dfde6cf\",\"label\":\"w-lnLi\",\"data\":0.8462610358696868},{\"uuid\":\"569ff181-8daa-42e4-8e90-d21188fd9f42\",\"label\":\"w-WCxt\",\"data\":0.39880772524949665},{\"uuid\":\"83376333-8d1a-47e0-be79-dba2c35f4222\",\"label\":\"w-UKsR\",\"data\":0.07396370986479524}],\"bias\":{\"uuid\":\"26989ef0-149b-406a-8931-8daf59e2521f\",\"label\":\"b-dMzy\",\"data\":-0.020826647213134656}}]},{\"neurons\":[{\"weights\":[{\"uuid\":\"417a3656-7e00-4cd8-b4de-f1f27f5c3394\",\"label\":\"w-oXTN\",\"data\":-0.5935180636439532},{\"uuid\":\"b20c3b3d-3e08-4d8e-9d8b-504f2421ee74\",\"label\":\"w-vbEn\",\"data\":-0.02483478285138913},{\"uuid\":\"2980f06c-1bfa-48e3-82b4-2ef215f9dfa5\",\"label\":\"w-sDaR\",\"data\":0.5267503685322347},{\"uuid\":\"55e64878-7432-4735-b4d4-494cc3fc2840\",\"label\":\"w-dPgL\",\"data\":-0.7360840806117643}],\"bias\":{\"uuid\":\"9f3384d9-f66f-4e0f-ad86-1122f2b48fcc\",\"label\":\"b-qtyX\",\"data\":0.6631648664486604}},{\"weights\":[{\"uuid\":\"004b4f22-6a7f-4d4e-8555-60c91e3d1e2a\",\"label\":\"w-faIj\",\"data\":0.05920513557112704},{\"uuid\":\"206cf548-fdfb-477d-bd58-03aaea3111f0\",\"label\":\"w-zPPR\",\"data\":0.6167016329015542},{\"uuid\":\"b58f1d5f-22b6-4c54-b077-a784d708225e\",\"label\":\"w-oXcr\",\"data\":-0.4281103307968058},{\"uuid\":\"94a52e9c-9750-4bb5-a463-fdc4c49e3f10\",\"label\":\"w-PqEk\",\"data\":-0.015538505374054873}],\"bias\":{\"uuid\":\"c3f55901-b34d-4b7f-bbb8-0c480362a4bc\",\"label\":\"b-KhWI\",\"data\":0.3499495296261822}},{\"weights\":[{\"uuid\":\"4aadd4c3-5b51-402c-b9ef-2237a84936a7\",\"label\":\"w-UoAs\",\"data\":0.4306733457097087},{\"uuid\":\"08ab1144-289f-4921-bd2b-bab1643e5b2b\",\"label\":\"w-zvfi\",\"data\":0.43290315934264645},{\"uuid\":\"188b1db4-9add-4734-879c-6754191d1fca\",\"label\":\"w-dIJb\",\"data\":1.2148152694389978},{\"uuid\":\"ff3e6050-b035-43d6-9863-e3c148bbfb24\",\"label\":\"w-fLqX\",\"data\":1.0015588721605535}],\"bias\":{\"uuid\":\"dbb29cb5-d311-4640-95bf-be45510b49e0\",\"label\":\"b-BQMY\",\"data\":-0.2798492033152966}},{\"weights\":[{\"uuid\":\"88bcf0ab-b11d-481f-8a2a-367817382b56\",\"label\":\"w-GVFk\",\"data\":-0.904381255042939},{\"uuid\":\"d8231755-16b0-4e08-9438-eeb4030726ec\",\"label\":\"w-XfvZ\",\"data\":0.7436924653998607},{\"uuid\":\"fa1c9cf4-de52-4983-8a2e-844a5ab05819\",\"label\":\"w-QluM\",\"data\":0.12537056954369105},{\"uuid\":\"4e4be4c0-eb99-440b-9349-d7abd7d091e8\",\"label\":\"w-DWQk\",\"data\":-0.8220082187353072}],\"bias\":{\"uuid\":\"7ad30add-edbf-4e6b-b725-6448c576a497\",\"label\":\"b-MDGo\",\"data\":-0.7167535944978365}}]},{\"neurons\":[{\"weights\":[{\"uuid\":\"eb427309-3e02-4fa4-9eca-a4aaa37135e2\",\"label\":\"w-NrRi\",\"data\":1.0528436734471773},{\"uuid\":\"bfabedce-0575-4630-b79b-800762169bfd\",\"label\":\"w-JYyD\",\"data\":0.5740347838184844},{\"uuid\":\"38cf1965-df0f-4126-8368-cceda2a17b30\",\"label\":\"w-GQIj\",\"data\":1.354121999397353},{\"uuid\":\"6a58aaaa-83c5-4769-97e0-2937e8a00c00\",\"label\":\"w-HLlS\",\"data\":0.7225452323588146}],\"bias\":{\"uuid\":\"e70308a9-849d-4d59-a4d9-0e141ac09423\",\"label\":\"b-nEgH\",\"data\":0.3714674363747061}}]}]}",
				[][]*Value{
					NewInputs([]float64{1.0, 3.0, -1.0}),
					NewInputs([]float64{3.0, -1.0, 0.5}),
					NewInputs([]float64{0.5, 1.0, 1.0}),
					NewInputs([]float64{1.0, 1.0, -1.0}),
				},
			},
			[][]float64{
				{0.8865051191267113},
				{-0.9087376153294349},
				{-0.8390334906148976},
				{0.8322463500127849},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mlp, err := FromJSON([]byte(tt.args.mlp))
			if err != nil {
				t.Errorf("FromJSON() = %v", err)
			}

			outputs, err := ForwardMLPMultipleInputs(mlp, tt.args.inputs)
			if err != nil {
				t.Errorf("FromJSON() = %v", err)
			}
			got := make([][]float64, len(tt.want))

			for i, subrange := range outputs {
				fmt.Printf("subrange i has length %v\n", len(subrange))
				got[i] = make([]float64, len(subrange))
				for j, value := range subrange {
					got[i][j] = value.Data
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("FromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

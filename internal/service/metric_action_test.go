package service

//type senderMock struct {
//	t   *testing.T
//	url string
//}
//
//func (s *senderMock) SendPost(ctx context.Context, url string, storage entity.Metric) {
//	if url != s.url {
//		assert.Fail(s.t, "url not equal")
//	}
//}
//
//func newSenderMock(t *testing.T, url string) *senderMock {
//	return &senderMock{
//		t:   t,
//		url: url,
//	}
//}
//
//func TestMetricAction(t *testing.T) {
//	type fields struct {
//		MemStorage *storage.MemStorage
//		sender     sender
//	}
//	type args struct {
//		ctx context.Context
//		url string
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		{
//			name: "success gauge metric",
//			args: args{
//				ctx: context.Background(),
//				url: "http://localhost:8080/update/gauge/Alloc/134",
//			},
//			fields: fields{
//				MemStorage: storage.NewMemStorage(),
//				sender:     sender2.NewSender(),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			m := &MetricAction{
//				MemStorage: tt.fields.MemStorage,
//				sender:     newSenderMock(t, tt.args.url),
//			}
//			m.SendMetric(tt.args.ctx, ":8080")
//		})
//	}
//}

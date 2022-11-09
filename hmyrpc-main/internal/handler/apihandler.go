package handler

import (
	"encoding/json"
	"net/http"

	"github.com/shenzhendev/hmyrpc/internal/logic"
	"github.com/shenzhendev/hmyrpc/internal/svc"
	"github.com/shenzhendev/hmyrpc/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ApiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		//if err := httpx.Parse(r, &req); err != nil {
		//	httpx.Error(w, err)
		//	return
		//}
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		l := logic.NewApiLogic(r.Context(), svcCtx)
		resp, err := l.HandleApiRequest(req)
		if err != nil {
			println("error processing", err.Error())
			httpx.OkJson(w, l.InvalidResp("2.0", req.ID))
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

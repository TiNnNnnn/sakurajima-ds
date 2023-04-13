package objects

import (
	"net/http"
)

func post(w http.ResponseWriter, r *http.Request) {
	// //获取对象名称
	// name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// //获取报头大小
	// size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }
	// //获取客户端上传的 object hash
	// hash := utils.GetHashFromHeader(r.Header)
	// if hash == "" {
	// 	log.Println("here is no object hash in digest header!")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }

}

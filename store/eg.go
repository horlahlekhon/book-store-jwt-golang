package store
//
//body, err := ioutil.ReadAll(r.Body)
//if err != nil {
//log.WithFields(logFields).Errorf("can't read request body: %s", err.Error())
//
//w.WriteHeader(http.StatusInternalServerError)
//fmt.Fprint(w, util.ParseError("server_error", fmt.Sprintf("can't read request body: %s", err.Error()), ""))
//
//return
//}
//
//err = r.Body.Close()
//if err != nil {
//log.WithFields(logFields).Errorf("can't close request: %s", err.Error())
//
//w.WriteHeader(http.StatusInternalServerError)
//fmt.Fprint(w, util.ParseError("server_error", fmt.Sprintf("can't close request: %s", err.Error()), ""))
//
//return
//}
//var posttedCheckbook PostedCheckbook
//err = json.Unmarshal(body, &postedCheckbook)
//if err != nil {
//log.WithFields(logFields).Errorf("can't parse request body:%s", err.Error())
//
//w.WriteHeader(http.StatusBadRequest)
//fmt.Fprint(w, util.ParseError("parse_error", fmt.Sprintf("can't parse request body: %s", err.Error()), ""))
//
//return
//}

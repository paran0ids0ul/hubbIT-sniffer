
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';

class myMacs extends WebComponent {
  
  
  @observable
  String cid;
  
  List<String> macs = toObservable(new List<String>());
  
  myMacs(){
      new Timer.periodic(new Duration(seconds:1), (Timer t) => makeRequest());
  }
  
  void makeRequest() {
    print("makerequers"+"");
    macs.clear();
    var path = "http://127.0.0.1:8080";
    var stuff = '?"site"="getAllMyMacs"';
    var json = '&"json"={"cid":"$cid"}';
    print(path+stuff+json);
    var request = HttpRequest.getString(path+stuff+json).then(requestComplete);
  }
  
  void requestComplete(String request) {
    var json = parse(request);
    print(json);
    //List<String> tmp = new List();
   // macs=json['mac'];
    for(var o in json['mac']){
      print(o);
      macs.add(o);
    }
  }
  
}

  


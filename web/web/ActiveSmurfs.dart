
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';

class ActiveSmurfs extends WebComponent {
  
  @observable
  
  List<String> smurfs = toObservable(new List<String>());
  
  ActiveSmurfs(){
  //  query('#button1').onClick.listen(makeRequest);
    
  // new Timer.periodic(new Duration(seconds:1), (Timer t) => makeRequest());
    
    
  }
  

  
  void makeRequest() {
    smurfs.clear();
    var path = "http://127.0.0.1:8080";
    var stuff = '?"site"="getSmurrfInTheHubb"';
    print(path);
    var request = HttpRequest.getString(path+stuff).then(requestComplete);
  }
  
  void requestComplete(String request) {
    var json = parse(request);
    List<String> tmp = new List();
    for(var o in json){
      smurfs.add(o["cid"]);
    }
  }
  
}

  


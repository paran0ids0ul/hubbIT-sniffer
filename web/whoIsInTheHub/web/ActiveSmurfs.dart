
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';

class ActiveSmurfs extends WebComponent {
  
  @observable
  
  List<String> smurfs = toObservable(new List<String>());
  
  ActiveSmurfs(){
    smurfs.add("Smurf1");
    smurfs.add("Smurf2");
    smurfs.add("Smurf3");
    smurfs.add("Smurf4");
    smurfs.add("Smurf5");
  //  query('#button1').onClick.listen(makeRequest);
    
   new Timer.periodic(new Duration(seconds:1), (Timer t) => makeRequest());
    
    
  }
  

  
  void makeRequest() {
    smurfs.clear();
    var path = "http://localhost:8080";
    var stuff = '?"site"="getSmurrfInTheHubb"';
    print(path);
    /*
    var httpRequest = new HttpRequest();
    httpRequest.open('GET', path);
    httpRequest.onLoadEnd.listen((e) => requestComplete(httpRequest));
    httpRequest.send('');*/
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

  


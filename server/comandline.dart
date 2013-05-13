import 'package:sqljocky/sqljocky.dart';
import 'package:sqljocky/utils.dart';
import 'dart:io';
import 'dart:uri';
import 'dart:json';

ConnectionPool pool;

ConnectionPool getPool(){
  if(pool==null)
    pool = new ConnectionPool(user:"root", password:"monraket", port:3306, db:"whoIsInTheHubb", host:"localhost");
  return pool;
}


main() {
  print("starting server");
  HttpServer.bind('129.16.180.213', 8080).then((server) {
    server.listen((HttpRequest request) {
      var get;
      String string2=request.queryParameters.toString();
      try{
        get=parse(string2);   
      } catch(e) {          
        print('error pars: $e');// No specified type, handles all
        print(request.queryParameters.toString());

        request.response.write("error json");
        request.response.write(request.queryParameters.toString());
      }
      if(get==null){
        print(get);
        request.response.write("get eamty");
        request.response.close();
        return;
      }
      print("here?");
      //print(get);
      String site=get['site'];
      if(site==null){
        request.response.write("site inte satt");
        request.response.close();
        return;
      }
      String json=get['json'];
      
      if(site==null){
        return;
      }
      
      print(site);
      switch (site) {
        case 'addMacs':
          addMacs(json,request);
          break;
        case 'topList':
          topList(json,request);
          break;
        case 'addMacAddress':
          addMacAddress(json,request);
          break;
        case 'getAllMyMacs':
          //http://129.16.180.213:8080/?%22site%22=%22getAllMyMacs%22&%22json%22={%22cid%22:%22erikax%22}
          getAllMyMacs(json,request);
          break;
        case 'getSmurrfInTheHubb':
          getSmurrfInTheHubb(json,request);
          break;
        default:
          request.response.write("felaktig site satt"); 
      }
     
      
    });
  });
}
void addMacs(json,request){
  var pool =getPool();
  //print(json);
  //print(json['mac']);
  for (String mac in json['mac']) {
   // print(s);
    pool.query('insert INTO macs (mac) VALUES ("$mac")').then((result) {
      print(mac);
    }).catchError((onError){
      print("dubblet");
    });
  } 
  request.response.write("macs added");
  request.response.close();
}
void topList(json,request){
  print("printing toplist");
  var pool =getPool();
  String stringReturn="";
  var list = new List();
  pool.query('select cid as cid,points as point from macadresses ORDER BY points DESC').then((result) {
    // request.response.write('Hello, world');
 //   return returnString="";
    for (var row in result) {
      String cid=row[0];
      String point=row[1];
      var user = {
       'cid' : cid,
       'point' : point
      };
      list.add(user);
      
      request.response.write(stringify(user));
    }
    request.response.close();
  });
    
}
//http://129.16.180.213:8080/?%22site%22=%22addMacAddress%22&%22json%22={%22mac%22:[%2200:21:63:b5:1f:ff%22]}
void addMacAddress(json,request){
  var pool =getPool();
  for (String mac in json['mac']) {
    pool.query('insert INTO macs (mac) VALUES ("$mac")').then((result) {
      print(mac);
    }).catchError((onError){
      print("dubblet");
    });
  } 
  request.response.write("macs added");
  request.response.close();
  
}

void getAllMyMacs(json,request){
  print("get all my macs not rdy yet");
  if(json['cid']==null){
    request.response.write("error nead to give cid");
    request.response.close();
    return;
  }
  var pool =getPool();
  String cid=json['cid'];
  String stringReturn="";
  var list = new List();
  String query='select mac from macToCid WHERE `cid`="$cid"';
  print("query: $query");
  pool.query(query).then((result) {
    for (var row in result) {
      print("row $row");
      String mac=row[0]; 
     // print("mac $mac");
      list.add(mac);
    }
    print(stringify(list));
    request.response.write(list);
    request.response.close();
  });
  

}
void getSmurrfInTheHubb(json,request){
  print("get all the smurrf in the hubb");
  var pool =getPool();
  String stringReturn="";
  var list = new List();
  pool.query('SELECT * FROM `macadresses` WHERE `timeInHubben`!=0').then((result) {
    for (var row in result) {
      String cid=row[0];
      String point=row[1];
      var user = {
            'cid' : cid
      };
      list.add(stringify(user));
    }
    request.response.write(list);
    request.response.close();
  });
}
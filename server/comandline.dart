import 'package:sqljocky/sqljocky.dart';
import 'package:sqljocky/utils.dart';
import 'dart:io';
import 'dart:json';

ConnectionPool pool;

ConnectionPool getPool(){
  if(pool==null)
    pool = new ConnectionPool(user:"root", password:"monraket", port:3306, db:"whoIsInTheHubb", host:"localhost");
  return pool;
}


main() {
  print("starting server");
  HttpServer.bind('129.16.186.148', 8080).then((server) {
    server.listen((HttpRequest request) {
      var get;
      try{
        get=parse(request.queryParameters.toString());
        //print(get);
        
      } catch(e) {                          // No specified type, handles all
        print('error pars');

        request.response.write("error json");
        request.response.write(request.queryParameters.toString());
       // request.response.close();
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
          TopList(json,request);
          break;
        case 'addMacAddress':
          addMacAddress(json,request);
          break;
        case 'getAllMyMacs':
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
void TopList(json,request){
  print("printing toplist");
  var pool =getPool();
  String stringReturn="";
  var list = new List();
  pool.query('select cid as cid,points as point from macadresses ORDER BY points DESC').then((result) {
    // request.response.write('Hello, world');
 //   return returnString="";
    for (var row in result) {
      for (String col in row) {
        //request.response.write(col);
      }
      
      //request.response.write("<br />");
      String cid=row[0];
      String point=row[1];
      //print("{$cid,$point}");
      var user = {
       'cid' : cid,
       'point' : point
      };
      list.add(user);
      
      request.response.write(stringify(user));
      
      //request.response.write(stringify(list));
    }
   // request.response.write(stringify(list));
   // request.response.write(stringify(result));
    request.response.close();
  });
    
}
void addMacAddress(json,request){
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

void getAllMyMacs(json,request){
  print("get all my macs not rdy yet");
  var pool =getPool();
  String stringReturn="";
  var list = new List();
  pool.query('select cid as cid,points as point from macadresses ORDER BY points DESC').then((result) {
    // request.response.write('Hello, world');
 //   return returnString="";
    for (var row in result) {
      //request.response.write("<br />");
      String cid=row[0];
      String point=row[1];
      //print("{$cid,$point}");
      var user = {
       'cid' : cid,
       'point' : point
      };
      list.add(user);
      
      request.response.write(stringify(user));
      
      //request.response.write(stringify(list));
    }
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
      list.add(user);
     // request.response.write(stringify(user));
    }
    request.response.write(list);
    request.response.close();
  });
}
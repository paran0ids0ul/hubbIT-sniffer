// Auto-generated from registerMacs.html.
// DO NOT EDIT.

library registerMacs;

import 'dart:html' as autogenerated;
import 'dart:svg' as autogenerated_svg;
import 'package:web_ui/web_ui.dart' as autogenerated;
import 'package:web_ui/observe/observable.dart' as __observe;
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';



class registerMacs extends WebComponent with Observable  {
  /** Autogenerated from the template. */

  /** CSS class constants. */
  static Map<String, String> _css = {};

  /** This field is deprecated, use getShadowRoot instead. */
  get _root => getShadowRoot("registerMacs");
  static final __shadowTemplate = new autogenerated.DocumentFragment.html('''
        <div>
          cid: <input type="text" placeholder="erikax">
          mac: <input type="text" placeholder="74:f0:6d:81:5f:cc">
          
          <button>add</button><br>
          </div>
      ''');
  autogenerated.ButtonElement __e15;
  autogenerated.InputElement __e13, __e14;
  autogenerated.Template __t;

  void created_autogenerated() {
    var __root = createShadowRoot("registerMacs");
    __t = new autogenerated.Template(__root);
    __root.nodes.add(__shadowTemplate.clone(true));
    __e13 = __root.nodes[1].nodes[1];
    __t.listen(__e13.onInput, ($event) { cid = __e13.value; });
    __t.oneWayBind(() => cid, (e) { if (__e13.value != e) __e13.value = e; }, false, false);
    __e14 = __root.nodes[1].nodes[3];
    __t.listen(__e14.onInput, ($event) { mac = __e14.value; });
    __t.oneWayBind(() => mac, (e) { if (__e14.value != e) __e14.value = e; }, false, false);
    __e15 = __root.nodes[1].nodes[5];
    __t.listen(__e15.onClick, ($event) { makeRequest(); });
    __t.create();
  }

  void inserted_autogenerated() {
    __t.insert();
  }

  void removed_autogenerated() {
    __t.remove();
    __t = __e13 = __e14 = __e15 = null;
  }

  /** Original code from the component. */

  
  
  String __$cid;
  String get cid {
    if (__observe.observeReads) {
      __observe.notifyRead(this, __observe.ChangeRecord.FIELD, 'cid');
    }
    return __$cid;
  }
  set cid(String value) {
    if (__observe.hasObservers(this)) {
      __observe.notifyChange(this, __observe.ChangeRecord.FIELD, 'cid',
          __$cid, value);
    }
    __$cid = value;
  }
  String mac;
 // List<String> macs = toObservable(new List<String>());
  
  registerMacs(){
     // new Timer.periodic(new Duration(seconds:1), (Timer t) => makeRequest());
  }
  
  void makeRequest() {
    print("makerequers"+"");
    var path = "http://127.0.0.1:8080";
    var stuff = '?"site"="registerMacAddress"';
    var json = '&"json"={"cid":"$cid","mac":"$mac"}';
    print(path+stuff+json);
    var request = HttpRequest.getString(path+stuff+json).then(requestComplete);
  }
  
  void requestComplete(String request) {
  
  }
  
}

  


//@ sourceMappingURL=registerMacs.dart.map
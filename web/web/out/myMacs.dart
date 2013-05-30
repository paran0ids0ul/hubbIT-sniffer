// Auto-generated from myMacs.html.
// DO NOT EDIT.

library myMacs;

import 'dart:html' as autogenerated;
import 'dart:svg' as autogenerated_svg;
import 'package:web_ui/web_ui.dart' as autogenerated;
import 'package:web_ui/observe/observable.dart' as __observe;
import 'package:web_ui/web_ui.dart';
import 'dart:html';
import 'dart:json';
import 'dart:async';



class myMacs extends WebComponent with Observable  {
  /** Autogenerated from the template. */

  /** CSS class constants. */
  static Map<String, String> _css = {};

  /** This field is deprecated, use getShadowRoot instead. */
  get _root => getShadowRoot("myMacs");
  static final __html1 = new autogenerated.BRElement(), __shadowTemplate = new autogenerated.DocumentFragment.html('''
        <div>
          cid: <input type="text" placeholder="erikax">
         <template></template>
        </div>
      ''');
  autogenerated.Element __e12;
  autogenerated.InputElement __e10;
  autogenerated.Template __t;

  void created_autogenerated() {
    var __root = createShadowRoot("myMacs");
    __t = new autogenerated.Template(__root);
    __root.nodes.add(__shadowTemplate.clone(true));
    __e10 = __root.nodes[1].nodes[1];
    __t.listen(__e10.onInput, ($event) { cid = __e10.value; });
    __t.oneWayBind(() => cid, (e) { if (__e10.value != e) __e10.value = e; }, false, false);
    __e12 = __root.nodes[1].nodes[3];
    __t.loop(__e12, () => macs, ($list, $index, __t) {
      var mac = $list[$index];
      var __binding11 = __t.contentBind(() => mac, false);
    __t.addAll([new autogenerated.Text(' '),
        __binding11,
        __html1.clone(true),
        new autogenerated.Text(' ')]);
    });
    __t.create();
  }

  void inserted_autogenerated() {
    __t.insert();
  }

  void removed_autogenerated() {
    __t.remove();
    __t = __e10 = __e12 = null;
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

  


//@ sourceMappingURL=myMacs.dart.map
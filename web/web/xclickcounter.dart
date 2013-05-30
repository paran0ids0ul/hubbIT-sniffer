import 'package:web_ui/web_ui.dart';

class CounterComponent extends WebComponent {
  int count = 0;

  CounterComponent(){
    smurfs.add("Smurf1");
    smurfs.add("Smurf2");
    smurfs.add("Smurf3");
    smurfs.add("Smurf4");
    smurfs.add("Smurf5");
  }
  
  
  void increment() {
    count++;
  }
  
 
  
  List<String> smurfs = new List<String>();
  

}

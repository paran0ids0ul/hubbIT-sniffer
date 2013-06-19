package whoisinthehubbserver

import grails.converters.JSON

class PointController {

    def topList() { 
		render(view: "create");
		def allUser = User.getAll();
		render allUser as JSON
	}
	def addCid(){
		String cidString = params.cid;
		String macAddr = params.mac;
		User user = User.find { cid == cidString };
		if(user==null){
			User u=new User(cid:cidString);
			u.save();
			def allUser = User.getAll();
			render allUser as JSON;
		}
		
		Mac macClass= new Mac(mac:macAddr);
		def allMac = Mac.getAll();
		render allMac as JSON;
	}
	def addTick(){
		
		
	}
	def getAllMyMacs(){
		String cidString = params.cid;
		User user = User.find { cid == cidString };
		if(user==null){
			render "cant find cid:  "+cidString;
		}else{
			render user.macs as JSON;
			
		}
	}
	
}

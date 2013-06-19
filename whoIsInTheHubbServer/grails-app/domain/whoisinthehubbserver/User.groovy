package whoisinthehubbserver

class User {
	String cid;
	int points=0;
	static hasMany = [ macs : Mac ]
    static constraints = {
    }
}

dataSource {
    pooled = true
    driverClassName = "org.h2.Driver"
    username = "sa"
    password = ""
}
hibernate {
    cache.use_second_level_cache = true
    cache.use_query_cache = false
    cache.region.factory_class = 'net.sf.ehcache.hibernate.EhCacheRegionFactory'
}
// environment specific settings
environments {
    development {
		dataSource {
			pooled = true
			dbCreate = "create-drop"
			url = "jdbc:mysql://localhost/whoIsInTheHubb"
			driverClassName = "com.mysql.jdbc.Driver"
			username = "root"
			password = "monraket"
		}
		/*dataSource {
            dbCreate = "create-drop" // one of 'create', 'create-drop', 'update', 'validate', ''
            url = "jdbc:h2:mem:devDb;MVCC=TRUE;LOCK_TIMEOUT=10000"
        }*/
    }
    test {
		
		dataSource {
			pooled = true
			dbCreate = "update"
			url = "jdbc:mysql://localhost/whoIsInTheHubb"
			driverClassName = "com.mysql.jdbc.Driver"
			username = "root"
			password = "monraket"
		}
		
       // dataSource {
       //     dbCreate = "update"
       //     url = "jdbc:h2:mem:testDb;MVCC=TRUE;LOCK_TIMEOUT=10000"
       // }
    }
    production {
		dataSource {
			pooled = true
			dbCreate = "update"
			url = "jdbc:mysql://localhost/whoIsInTheHubb"
			driverClassName = "com.mysql.jdbc.Driver"
			username = "root"
			password = "monraket"
		}
		
		
		/*
        dataSource {
            dbCreate = "update"
            url = "jdbc:h2:prodDb;MVCC=TRUE;LOCK_TIMEOUT=10000"
            pooled = true
            properties {
               maxActive = -1
               minEvictableIdleTimeMillis=1800000
               timeBetweenEvictionRunsMillis=1800000
               numTestsPerEvictionRun=3
               testOnBorrow=true
               testWhileIdle=true
               testOnReturn=true
               validationQuery="SELECT 1"
            }
        }*/
    }
}

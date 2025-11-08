resource "persistent_buckets" "example" {
  bucket_capacity = 100
  maximum_buckets = 2
  items = {
    item-1 = {
      weight = 50
      item   = "some string data here"
    }
    item-2 = {
      weight = 25
      item   = null
    }
    item-3 = {
      weight = 10
    }
    item-4 = {
      weight = 50
    }
    item-5 = {
      weight = 10
    }
  }
}

# Result: 
# persistent_buckets.example.buckets = [
#   {  
#     item-1 = {
#  		  weight = 50
# 		  item   = "some string data here"
# 	  } 
#     item-2 = {
#   	  weight = 25
# 		  item   = null
#     } 
#     item-3 = {
# 	   	weight = 10
#	    } 
#     item-5 = {
#	  	  weight = 10
# 	  } 
#  },
#  {
#    item-4 = {
#		   weight = 50
#	   } 
#  } 
# ]

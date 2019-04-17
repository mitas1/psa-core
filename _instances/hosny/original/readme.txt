#
#
# The associated files are created by Manar Hosny to be used as test data for the single vehicle pickup and delivery problem with time windows.
# Please contact me if you need any help with the data: manar_hosny@hotmail.com
# The file name indicates the number of requests in the file. for example PDP_100.txt is a file with 100 requests, where each request is 
# a pair consisting of a pickup and an associated delivery.

# The file is organized as follows:

Vehicle Capacity
Request_number		X_coordinate		Y_coordinate		Demand		Early_Time_Window		Late_Time_Window

# Request #0 is the depot. Its demand is 0, and its time window interval is very large, i.e. no restriction  on the time window of the depot.
# The first half of the requests are all pickup requests, while the second half are all delivery requests, such that if the number of requests is n, 
# pickup request  i has an associated delivery request  i+n.
# For example, if the number of requests is 100, pickup request 1 is associated with delivery request 101, pickup request 2 is associated with 
# delivery request 102, ...etc.
# The demand of a delivery is the same as the demand of the pickup but with a negative sign. 
# The total number of locations in the file including the depot is 2n+1.
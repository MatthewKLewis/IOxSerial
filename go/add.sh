echo "Adding New..." 
&& ioxclient docker package postposter ./conf 
&& cd conf
&& ioxclient application run portposter package.tar --payload activation.json
&& cd ..
&& echo "Listing Active..." 
&& ioxclient application list
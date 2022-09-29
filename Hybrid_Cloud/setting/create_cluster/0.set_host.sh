echo -n "Your Hostname? "
read hostname

echo `hostname -i`      $hostname >> /etc/hosts
hostnamectl set-hostname $hostname

echo "Please Reconnect shell to Change Hostname!!"

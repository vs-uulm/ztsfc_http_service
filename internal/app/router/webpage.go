package router

var webpage string = `
<!DOCTYPE html>
<html>
<head>
<title>bwNET2020+</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap" rel="stylesheet">
<style>
* {margin: 0; padding: 0;}
body {font-family: 'Roboto', sans-serif;}
h1 {text-align: center; font-size: 50px; margin-bottom: 30px; font-weight: 500;}
h2 {text-align: center; font-size: 25px; margin-bottom: 20px; font-weight: 500;} p {font-size: 18px; margin-bottom: 20px; line-height: 1.5;} p:last-child{margin-bottom: 0;}  .grey {background-color:#F0F2F2;} .black {background-color:black;} .black * {color: #fff;} .yellow {background-color: #FFDE0D;}  .cewrap {padding: 50px 25px; } .cewrap > div {max-width: 1200px;position: relative; margin:0 auto;text-align: center;} </style> </head> 
<body>
<main>
<div class="cewrap yellow">
<div>
<h1>bwNET2020+</h1>
<h2>Research and innovative services for a flexible network in  <br> Baden-Württemberg</h2> <p> <br>In order to provide end users a familiar quality of experience the underlying network has to be equipped with new innovative technologies and solutions. The state project bwNet2020+ aims to leverage new innovative approaches in order to support the expansion and modernization of the BelWü network and university networks.
</p>
</div>
</div>

<div class="cewrap grey">
<div>
<h2>Self-driving Networks</h2> <p>The progressive digitalization that is emerging demands for secure and robust high-performance communication networks in this day and age. Those networks have long become an essential infrastructure equally needed by society and economy. Increasing data rates (up to 1 Tbit/s) and advances in telecommunications and cellular networks like 5G offer new possibilities to support applications. For example, service function chaining (SFC) allows services like load balancers or firewalls to be connected to a chain through which data packets should travel before reaching the final target. Service function chaining together with technologies like software-defined networks (SDN), programmable data planes (P4) and network function virtualization (NFV) lay the foundation for self-driving networks. By utilizing policies and monitoring data enhanced by machine learning, self-driving networks are able to run, adapt and defend themselves.
</p>
<p>bwNET2020+ aims to leverage the just mentioned technologies in order to advance current implementations of the BelWü network as well as university networks. The goal is to present viable options to transition those solutions gradually towards full adoption. Multiple use cases will be developed together by network administrators and network scientists.</p>
</div>
</div>


<div class="cewrap yellow">
<div>
<h2>Related Projects</h2>
<p>The project partners are able to build up on extensive experience and successful cooperation from the previous projects bwNET100G+, bwNET100G+ Extension and bwNetFlow. The challenges to the BelWü network were successfully accompanied by the implementation of 100 Gbit/s communication links at the universities in Baden-Württemberg. 
</p>
</div>
</div>
</main>
<footer class="cewrap black">
<div>
<h2> Contact</h2>
<p>
Ulm University<br>
Kommunikations- und Informationszentrum<br>
Albert-Einstein-Allee 11<br>
89081 Ulm - Germany</p>
</div>
</footer>
</body>
</html>

    ` 


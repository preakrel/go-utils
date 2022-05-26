<?php

parse_str("id=187923&poc[][0]=a&poc[][1]=b&poc[][2]=c&poc[]=d",$arr);
print_r(json_encode($arr));
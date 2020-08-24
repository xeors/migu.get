<?
    function Http($url,$data=null,$header=null,$rh=0,$nb=0){
        $ch=curl_init();
        curl_setopt($ch,CURLOPT_URL,$url);
        curl_setopt($ch,CURLOPT_HEADER,$rh);
        curl_setopt($ch,CURLOPT_NOBODY,$nb);
        curl_setopt($ch,CURLOPT_RETURNTRANSFER,1);
        $header==null?:curl_setopt($ch,CURLOPT_HTTPHEADER,$header);
        $data==null?:(curl_setopt($ch,CURLOPT_POST,1)&&curl_setopt($ch,CURLOPT_POSTFIELDS,$data));
        $rdata=curl_exec($ch);
        curl_close($ch);
        return $rdata;
    }
    function tr($data){
        if(preg_match('/.*?(curl).*?/',$_SERVER['HTTP_USER_AGENT'])!=0){
            return $data;
        }
        return '<br><br><br><center><textarea style="width:75vw;height:75vh">'.$data.'</textarea></center>';
    }
    if(isset($_GET['s'])){
        $sS='{"song":1,"album":0,"singer":0,"tagSong":1,"mvSong":0,"songlist":0,"bestShow":1,"lyricSong":0,"concert":0,"periodical":0,"ticket":0,"bit24":0,"verticalVideoTone":0}';
        $url='http://jadeite.migu.cn:7090/music_search/v2/search/searchAll?isCopyright=1&isCorrect=1&pageIndex=1&pageSize='.(isset($_GET['l'])?$_GET['l']:10).'&searchSwitch='.urlencode($sS).'&text='.urlencode($_GET['s']);
        $header=array(
            'timeStamp: 1596391260',
            'sign: 98d131a54422139907d45f7f204ecf72'
        );
        $rsp=Http($url,null,$header);
        if(!isset($_GET['r'])){
            $r=json_decode($rsp,1);
            $R=$r['songResultData']['result'];
            for($i=0;$i<count($R);$i++){
                $Sn=$R[$i]['songName'];
                $Sg=$R[$i]['singer'];
                $Al=$R[$i]['album'];
                //$Ds=$R[$i]['songDescription'];
                $Su=$R[$i]['newRateFormats'];
                echo "$Sn - $Sg - $Al <br>";
                for($j=0;$j<count($Su);$j++){
                    $Ft=$Su[$j]['formatType'];
                    if($Ft=='SQ' || $Ft=='ZQ'){
                        $u=preg_replace('/ftp:\/\/.*?\//','http://freetyst.nf.migu.cn/',$Su[$j]['androidUrl']);
                        echo '<a target="_blank" href="'.$u.'">'.$Ft.'</a> / ';
                    }else{
                        $u=preg_replace('/ftp:\/\/.*?\//','http://freetyst.nf.migu.cn/',$Su[$j]['url']);
                        echo '<a target="_blank" href="'.$u.'">'.$Ft.'</a> / ';
                    }
                }
                echo '<br><br>';
            }            
        }else{
            echo $rsp;
        }
    }
?>
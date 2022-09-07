import it.unisa.dia.gas.jpbc.Element;
import it.unisa.dia.gas.jpbc.Pairing;
import it.unisa.dia.gas.plaf.jpbc.pairing.PairingFactory;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Base64;
import java.util.Properties;

public class sas {
    //初始化选取曲线G1、G2以及生成元g1和g2
    public static void setup(String pairingParametersFileName,String mpkFileName){
        Pairing bp = PairingFactory.getPairing(pairingParametersFileName);

        Element g1 = bp.getG1().newRandomElement().getImmutable();
        Element g2 = bp.getG2().newRandomElement().getImmutable();

        Properties mpkProp = new Properties();
        mpkProp.setProperty("g1", Base64.getEncoder().encodeToString(g1.toBytes()));
        mpkProp.setProperty("g2", Base64.getEncoder().encodeToString(g2.toBytes()));
        storeProToFile(mpkProp,mpkFileName);

    }
    //the first node generate signature
    public static void jiami1(String pairingParametersFileName,String mpkFileName,String skFileName, String message,String sasFileName){
        Pairing bp = PairingFactory.getPairing(pairingParametersFileName);
        Properties mpkProp = loadPropFromFile(mpkFileName);
        String g2String = mpkProp.getProperty("g2");
        Element g2 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(g2String)).getImmutable();

        //the first node generate public key and private key
        Element x = bp.getZr().newRandomElement().getImmutable();
        Element v = g2.powZn(x);
        Properties skProp = new Properties();
        skProp.setProperty("x",Base64.getEncoder().encodeToString(x.toBytes()));//current node private key
        storeProToFile(skProp,skFileName);



        String c = "1" + message;//define e(∅, g2)=1in the c , the previous pubilc key is null,the message is null when the node is the first one.
        byte[] m_hash=Integer.toString(c.hashCode()).getBytes();//c的 hash chain value
        Element h = bp.getG1().newElementFromHash(m_hash, 0, m_hash.length).getImmutable();
        Element deta = h.powZn(x);//get the first node self-signature,and the self-signature=Aggregate signature because this node is the first node

        Properties sasProp = new Properties();//out put Aggregate signature inlcude(deta,pk,c,s)
        sasProp.setProperty("deta",Base64.getEncoder().encodeToString(deta.toBytes()));//Aggregate signature
        sasProp.setProperty("pk",Base64.getEncoder().encodeToString(v.toBytes()));//this node public key
        sasProp.setProperty("c",Base64.getEncoder().encodeToString(m_hash));//ciphertext
        sasProp.setProperty("s",Base64.getEncoder().encodeToString(deta.toBytes()));//self-signature
        storeProToFile(sasProp,sasFileName);
    }

    public static void jiami2(String pairingParametersFileName,String mpkFileName,String skFileName, String message,String sas0FileName,String sas1FileName)  {
        Pairing bp = PairingFactory.getPairing(pairingParametersFileName);
        Properties mpkProp = loadPropFromFile(mpkFileName);
        String g2String = mpkProp.getProperty("g2");
        Element g2 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(g2String)).getImmutable();

        Properties sasProp = loadPropFromFile(sas0FileName);
        String detaString = sasProp.getProperty("deta");
        Element detaPre = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(detaString)).getImmutable();

        String pkString = sasProp.getProperty("pk");
        Element pkPre = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(pkString)).getImmutable();

        String cString = sasProp.getProperty("c");
        byte[] cPre_hash = Base64.getDecoder().decode(cString);
        Element cPre = bp.getG1().newElementFromHash(cPre_hash, 0, cPre_hash.length).getImmutable();

        String sString = sasProp.getProperty("s");
        Element sPre = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(sString)).getImmutable();

        //Verify the previous signature
        if(bp.pairing(sPre,g2).isEqual(bp.pairing(cPre,pkPre))){
//            System.out.println(bp.pairing(sPre,g2).toString());
//            System.out.println(bp.pairing(cPre,pkPre).toString());
            System.out.println("the previous signature is valid");

            Element x = bp.getZr().newRandomElement().getImmutable();//generate current node private key
            Element v = g2.powZn(x);//generate current node public key
            Properties skProp = new Properties();
            skProp.setProperty("x",Base64.getEncoder().encodeToString(x.toBytes()));
            storeProToFile(skProp,skFileName);


            String c1 = bp.pairing(detaPre, g2).getImmutable().toString();
            String c2 =message;
            String c3=pkString;
            String c4=Arrays.toString(cPre_hash);
            String m=c1+c2+c3+c4;//generate current node ciphertext c
//            System.out.println("m:"+m);
            byte[] c_hash = Integer.toString(m.hashCode()).getBytes();
            Element c = bp.getG1().newElementFromHash(c_hash, 0, c_hash.length).getImmutable();
//            System.out.println("第一个c："+c);
            Element s = c.powZn(x).getImmutable();//generate current node self-signature
//            System.out.println("s:"+s);
            Element deta = detaPre.mul(s).getImmutable();//generate current node Aggregate signature

            sasProp.setProperty("deta",Base64.getEncoder().encodeToString(deta.toBytes()));
            sasProp.setProperty("pk",Base64.getEncoder().encodeToString(v.toBytes()));
            sasProp.setProperty("c",Base64.getEncoder().encodeToString(c_hash));
            sasProp.setProperty("s",Base64.getEncoder().encodeToString(s.toBytes()));
            storeProToFile(sasProp,sas1FileName);

        }else {
            System.out.println("the previous signature is invalid");
        }

    }


    public static Properties loadPropFromFile(String fileName) {
        Properties prop=new Properties();
        try(FileInputStream in =new FileInputStream(fileName)){
            prop.load(in);
        }catch (Exception e){
            e.printStackTrace();
            System.out.println(fileName+"load failed!");
            System.exit(-1);
        }
        return prop;
    }

    public static void storeProToFile(Properties prop, String fileName) {
        try(FileOutputStream out =new FileOutputStream(fileName)){
            prop.store(out,null);
        }catch (Exception e){
            e.printStackTrace();
            System.out.println(fileName+"save failed!");
            System.exit(-1);
        }
    }

    public static void getMessage (String filePath) {
        try {
            File file = new File(filePath);
            FileInputStream readIn = new FileInputStream(file);
            InputStreamReader read = new InputStreamReader(readIn, "utf-8");
            BufferedReader bufferedReader = new BufferedReader(read);
            String oneLine= null;
            while((oneLine= bufferedReader.readLine()) != null){
                System.out.println(oneLine);

            }
            read.close();
        } catch (Exception e) {
            System.out.println("read file error!");
            e.printStackTrace();
        }
    }

    public static void main(String[] args) throws Exception {
        String dir ="data2/";
        String peer1 ="peer1/";
        String peer2 ="peer2/";
        String peer3 ="peer3/";
        String peer4 ="peer4/";
        String peer5 ="peer5/";
        String peer6 ="peer6/";
        String peer7 ="peer7/";
        String peer8 ="peer8/";
        String peer9 ="peer9/";
        String peer10 ="peer10/";
        String peer11 ="peer11/";
        String peer12 ="peer12/";
        String peer13 ="peer13/";
        String peer14 ="peer14/";
        String peer15 ="peer15/";
        String peer16 ="peer16/";
        String peer17 ="peer17/";
        String peer18 ="peer18/";
        String peer19 ="peer19/";
        String peer20 ="peer20/";



        String pairingParametersFileName="a.properties";
        String mpkFileName=peer1+"mpk.properties";
        String skFileName_peer1 = peer1 +"sk.properties";//存储第一个节点私钥
        String sasFileName_peer1 = peer1 +"sas.properties";//存储第一个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer1 = peer1 +"message.properties";//存储第一个节点的消息

        String skFileName_peer2 = peer2 +"sk.properties";//存储第二个节点私钥
        String sasFileName_peer2 = peer2 +"sas.properties";//存储第二个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer2 = peer2 +"message.properties";//存储第一个节点的消息

        String skFileName_peer3 = peer3 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer3 = peer3 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer3 = peer3 +"message.properties";//存储第一个节点的消息

        String skFileName_peer4 = peer4 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer4 = peer4 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer4 = peer4 +"message.properties";//存储第一个节点的消息

        String skFileName_peer5 = peer5 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer5 = peer5 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer5 = peer5 +"message.properties";//存储第一个节点的消息

        String skFileName_peer6 = peer6 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer6 = peer6 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer6 = peer6 +"message.properties";//存储第一个节点的消息

        String skFileName_peer7 = peer7 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer7 = peer7 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer7 = peer7 +"message.properties";//存储第一个节点的消息

        String skFileName_peer8 = peer8 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer8 = peer8 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer8 = peer8 +"message.properties";//存储第一个节点的消息

        String skFileName_peer9 = peer9 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer9 = peer9 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer9 = peer9 +"message.properties";//存储第一个节点的消息

        String skFileName_peer10 = peer10 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer10 = peer10 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer10 = peer10 +"message.properties";//存储第一个节点的消息

        String skFileName_peer11 = peer11 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer11 = peer11 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer11 = peer11 +"message.properties";//存储第一个节点的消息

        String skFileName_peer12 = peer12 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer12 = peer12 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer12 = peer12 +"message.properties";//存储第一个节点的消息

        String skFileName_peer13 = peer13 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer13 = peer13 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer13 = peer13 +"message.properties";//存储第一个节点的消息

        String skFileName_peer14 = peer14 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer14 = peer14 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer14 = peer14 +"message.properties";//存储第一个节点的消息

        String skFileName_peer15 = peer15 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer15 = peer15 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer15 = peer15 +"message.properties";//存储第一个节点的消息

        String skFileName_peer16 = peer16 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer16 = peer16 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer16 = peer16 +"message.properties";//存储第一个节点的消息

        String skFileName_peer17 = peer17 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer17 = peer17 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer17 = peer17 +"message.properties";//存储第一个节点的消息

        String skFileName_peer18 = peer18 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer18 = peer18 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer18 = peer18 +"message.properties";//存储第一个节点的消息

        String skFileName_peer19 = peer19 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer19 = peer19 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer19 = peer19 +"message.properties";//存储第一个节点的消息

        String skFileName_peer20 = peer20 +"sk.properties";//存储第三个节点私钥
        String sasFileName_peer20 = peer20 +"sas.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名
        String message_peer20 = peer20 +"message.properties";//存储第一个节点的消息



//        String mpkFileName = dir +"mpk.properties";//存储生成元g1和g2
//        String skFileName = dir +"sk_0.properties";//存储第一个节点私钥
//        String sasFileName = dir +"sas_0.properties";//存储第一个节点的聚合签名，公钥，消息，未聚合的签名
//        String skFileName_1 = dir +"sk_1.properties";//存储第二个节点私钥
//        String sasFileName_1 = dir +"sas_1.properties";//存储第二个节点的聚合签名，公钥，消息，未聚合的签名
//        String skFileName_2 = dir +"sk_2.properties";//存储第三个节点私钥
//        String sasFileName_2 = dir +"sas_2.properties";//存储第三个节点的聚合签名，公钥，消息，未聚合的签名

//        String s1="第一条消息";
//        String s2="第二条消息";
//        String s3="第三条消息";


        String s1 = new String(Files.readAllBytes(Paths.get("peer1/message.txt")));
//        System.out.println(s1);
        String s2 = new String(Files.readAllBytes(Paths.get("peer2/message.txt")));
//        System.out.println(s2);
        String s3 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
//        System.out.println(s3);
        System.out.println("===================");
        String s4 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s5 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s6 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s7 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s8 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s9 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s10 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s11 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s12 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s13 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s14 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s15 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s16 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s17 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s18 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s19 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));
        String s20 = new String(Files.readAllBytes(Paths.get("peer3/message.txt")));



//        Properties messageProp = new Properties();
//        messageProp.setProperty("message", Base64.getEncoder().encodeToString(s1.getBytes()));
////        mpkProp.setProperty("g2", Base64.getEncoder().encodeToString(g2.toBytes()));
//        storeProToFile(messageProp,message_peer1);
//
//
//        Properties messageProp = loadPropFromFile("message_peer1");
//        String message = messageProp.getProperty("message");
////        Base64.getDecoder().decode(message);
//        System.out.println(message);
        long start = System.currentTimeMillis();
        //加密过程
        setup(pairingParametersFileName,mpkFileName);
        System.out.println("The first node:");
        jiami1(pairingParametersFileName,mpkFileName,skFileName_peer1,s1,sasFileName_peer1);
        System.out.println("The second node:");//需要拿到上一个节点的签名输出文件和生成元g2，函数生成一个私钥文件，一个有序聚合签名文件（包括聚合签名，当前节点公钥，聚合消息哈希和未聚合签名）
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer2,s2,sasFileName_peer1,sasFileName_peer2);
        System.out.println("The third node:");
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer3,s3,sasFileName_peer2,sasFileName_peer3);
//        4
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer4,s4,sasFileName_peer3,sasFileName_peer4);

//        5
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer5,s5,sasFileName_peer4,sasFileName_peer5);

//        6
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer6,s6,sasFileName_peer5,sasFileName_peer6);

//        7
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer7,s7,sasFileName_peer6,sasFileName_peer7);

//        8
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer8,s8,sasFileName_peer7,sasFileName_peer8);

//        9
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer9,s9,sasFileName_peer8,sasFileName_peer9);

//        10
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer10,s10,sasFileName_peer9,sasFileName_peer10);
//        11
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer11,s11,sasFileName_peer10,sasFileName_peer11);
//        12
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer12,s12,sasFileName_peer11,sasFileName_peer12);
//        13
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer13,s13,sasFileName_peer12,sasFileName_peer13);
//        14
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer14,s14,sasFileName_peer13,sasFileName_peer14);
//        15
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer15,s15,sasFileName_peer14,sasFileName_peer15);
//        16
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer16,s16,sasFileName_peer15,sasFileName_peer16);
//        17
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer17,s17,sasFileName_peer16,sasFileName_peer17);
//        18
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer18,s18,sasFileName_peer17,sasFileName_peer18);
//        19
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer19,s19,sasFileName_peer18,sasFileName_peer19);
//        20
        jiami2(pairingParametersFileName,mpkFileName,skFileName_peer20,s20,sasFileName_peer19,sasFileName_peer20);



        long end = System.currentTimeMillis();

//        System.out.println("加密耗时："+(end-start));
        //解密过程


        Pairing bp = PairingFactory.getPairing(pairingParametersFileName);
        String c1="1"+s1;//pk0和c0为空集，e(H（null）,pk)=1,

        byte[] c1_hash = Integer.toString(c1.hashCode()).getBytes();
        Element Hc1 = bp.getG1().newElementFromHash(c1_hash, 0, c1_hash.length).getImmutable();//第一个Hc1

        Properties sas1Prop = loadPropFromFile(sasFileName_peer1);
        String pk1String = sas1Prop.getProperty("pk");
        Element pk1 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(pk1String)).getImmutable();

        Properties sas2Prop = loadPropFromFile(sasFileName_peer2);
        String pk2String = sas2Prop.getProperty("pk");
        Element pk2 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(pk2String)).getImmutable();

        Properties sas3Prop = loadPropFromFile(sasFileName_peer3);
        String pk3String = sas3Prop.getProperty("pk");
        Element pk3 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(pk3String)).getImmutable();

        String c2 = bp.pairing(Hc1, pk1).toString() + s2 + pk1String + Arrays.toString(c1_hash);
//        System.out.println("c2"+c2);
        byte[] c2_hash = Integer.toString(c2.hashCode()).getBytes();
        Element Hc2 = bp.getG1().newElementFromHash(c2_hash, 0, c2_hash.length).getImmutable();
//        System.out.println("第二个c:"+Hc2);

        String c3 = bp.pairing(Hc1, pk1).getImmutable().mul(bp.pairing(Hc2, pk2)).toString() + s3 + pk2String + Arrays.toString(c2_hash);
        byte[] c3_hash = Integer.toString(c3.hashCode()).getBytes();
        Element Hc3 = bp.getG1().newElementFromHash(c3_hash, 0, c3_hash.length).getImmutable();


        Properties mpkProp = loadPropFromFile(mpkFileName);
        String g2String = mpkProp.getProperty("g2");//获取生成元
        Element g2 = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(g2String)).getImmutable();
//        Properties sas3Prop = loadPropFromFile(sasFileName_2);
        String detaString = sas3Prop.getProperty("deta");//获取到第二个节点为止的聚合签名
        Element deta = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(detaString)).getImmutable();
        String sString = sas3Prop.getProperty("s");//获取到第二个节点的未聚合的签名
        Element s = bp.getG1().newElementFromBytes(Base64.getDecoder().decode(sString)).getImmutable();



        Element right = bp.pairing(Hc1, pk1).getImmutable().mul(bp.pairing(Hc2, pk2).mul(bp.pairing(Hc3,pk3))).getImmutable();
        Element left = bp.pairing(deta, g2).getImmutable();
        System.out.println("left:"+left);
        System.out.println("right:"+right);
//        System.out.println(bp.pairing(Hc1, g2_x1));
//        System.out.println(bp.pairing(Hc2, g2_x2));
//        System.out.println(bp.pairing(s,g2));
        if (left.isEqual(right)){
            System.out.println("SAS is valid");
        }else {
            System.out.println("SAS is invalid");
        }

    }
}



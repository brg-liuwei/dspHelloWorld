package aes

/*
#cgo CFLAGS: -I/usr/local/java/jdk1.7.0_67/include/ -I/usr/local/java/jdk1.7.0_67/include/linux/
#cgo LDFLAGS: -L/usr/local/java/jdk1.7.0_67/jre/lib/amd64/server/ -ljvm

#include <jni.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
#include <wchar.h>
#include <sys/types.h>

static JavaVM *jvm;
static JNIEnv *env;

void JVM_Init()
{
    JavaVMInitArgs vm_args;
    JavaVMOption options[1];

    vm_args.version = JNI_VERSION_1_2;
    vm_args.ignoreUnrecognized = JNI_TRUE;
    vm_args.nOptions = 0;

    char classpath[1024] = "-Djava.class.path=";
    char *env_classpath = getenv("CLASSPATH");

    if (env_classpath) {
        options[0].optionString = strcat(classpath, env_classpath);
        vm_args.nOptions++;
    }

    if (vm_args.nOptions > 0) {
        vm_args.options = options;
    }

    // 创建java虚拟机
    jint res = JNI_CreateJavaVM(&jvm, (void **)&env, &vm_args);
    if (res < 0) {
        printf("Create Java VM error, code = %d/n", res);
        exit(-1);
    }
}

void JVM_Destroy()
{
    (*jvm)->DestroyJavaVM(jvm);
    env = NULL;
    jvm = NULL;
}

char *jstring2ch(jstring jstr)
{
    char* rtn = NULL;
    jclass clsstring = (*env)->FindClass(env, "java/lang/String");
    jstring strencode = (*env)->NewStringUTF(env, "utf-8");
    jmethodID mid = (*env)->GetMethodID(env, clsstring, "getBytes", "(Ljava/lang/String;)[B");
    jbyteArray barr= (jbyteArray)(*env)->CallObjectMethod(env, jstr, mid, strencode);
    jsize alen = (*env)->GetArrayLength(env, barr);
    jbyte* ba = (*env)->GetByteArrayElements(env, barr, JNI_FALSE);
    if (alen > 0) {
    	rtn = (char*)malloc(alen + 1);
    	memcpy(rtn, ba, alen);
    	rtn[alen] = 0;
    }
    (*env)->ReleaseByteArrayElements(env, barr, ba, 0);
    return rtn;
}

jstring ch2jstring(const char* pat)
{
    // fprintf(stderr, "<<< ch2jstring-1\n");
    // jclass strClass = (*env)->FindClass(env, "java/lang/String");
    // fprintf(stderr, "<<< ch2jstring-2 strClass: %p\n", strClass);
    // jmethodID ctorID = (*env)->GetMethodID(env, strClass, "<init>", "([BLjava/lang/String;)V");
    // fprintf(stderr, "<<< ch2jstring-3 ctorId: %p\n", ctorID);
    // jbyteArray bytes = (*env)->NewByteArray(env, strlen(pat));
    // fprintf(stderr, "<<< ch2jstring-4\n");
    // (*env)->SetByteArrayRegion(env, bytes, 0, strlen(pat), (jbyte*)pat);
    // fprintf(stderr, "<<< ch2jstring-5\n");
    // jstring encoding = (*env)->NewStringUTF(env, "utf-8");
    // fprintf(stderr, "<<< ch2jstring-6\n");
    // return (jstring)(*env)->NewObject(env, strClass, ctorID, bytes, encoding);

    //jclass strClass = (*env)->FindClass(env, "java/lang/String");
    jstring str = (*env)->NewStringUTF(env, pat);
    return str;
}

static jclass getPriceEncryptHelper()
{
    static jclass helper = NULL;

    if (helper) {
        fprintf(stderr, "~~~~~~~~find helper: %p\n", helper);
        return helper;
    }

    fprintf(stderr, "find class helper: env = %p\n", *env);
    helper = (*env)->FindClass(env, "PriceEncryptHelper");
    fprintf(stderr, "find class helper after\n");
    if (!helper) {
        fprintf(stderr, "cannot find PriceEncryptHelper, exit\n");
        exit(-1);
    }
    fprintf(stderr, "find class helper = %p\n", helper);
    return helper;
}

static jclass getPriceEncryptInfo()
{
    static jclass info = NULL;

    if (info) {
        return info;
    }

    info = (*env)->FindClass(env, "PriceEncryptInfo");
    if (!info) {
        printf("class PriceEncryptInfo cannot be found\n");
        exit(-1);
    }
    return info;
}

static jmethodID getDecryptedPriceMethod(jclass helper)
{
    static jmethodID mid = NULL;

    if (mid) {
        fprintf(stderr, "~~~~ found static mid %p\n", mid);
        return mid;
    }

    mid = (*env)->GetStaticMethodID(env, helper, "getDecryptedPrice", "(Ljava/lang/String;Ljava/lang/String;)LPriceEncryptInfo;");
    if (!mid) {
        fprintf(stderr, "cannot find method getDecryptedPriceMethod\n");
        exit(-1);
    }
    fprintf(stderr, "find jmethodId mid: %p\n", mid);
    return mid;
}

static jfieldID getPriceFieldId(jclass info)
{
    static jfieldID priceId = NULL;

    if (priceId) {
        return priceId;
    }

    priceId = (*env)->GetFieldID(env, info, "price", "J");
    if (!priceId) {
        printf("cannot found field price\n");
        exit(-1);
    }
    return priceId;
}

long getDecryptedPrice(const char *ePrice, const char *key)
{
    fprintf(stderr, "--------------> aes: getDecryptedPrice invoke\n");
    jclass helper = getPriceEncryptHelper();
    fprintf(stderr, "--------------> aes: getHelper: %p\n", (void *)helper);
    jmethodID mid = getDecryptedPriceMethod(helper);
    fprintf(stderr, "==============> aes: get mid: %p\n", (void *)mid);
    jstring encryptedPrice = ch2jstring(ePrice);
    fprintf(stderr, "##############> aes: get jprice: %p\n", (void *)encryptedPrice);
    jstring jkey = ch2jstring(key);
    fprintf(stderr, "--------------> aes: get jkey: %p\n", (void *)jkey);

    jobjectArray args = (*env)->NewObjectArray(env, 2, (*env)->FindClass(env, "java/lang/String"), NULL);
    (*env)->SetObjectArrayElement(env, args, 0, encryptedPrice);
    (*env)->SetObjectArrayElement(env, args, 1, jkey);

    fprintf(stderr, "set price and key ok\n");
    //jobject info = (*env)->CallStaticObjectMethod(env, helper, mid, encryptedPrice, jkey);
    jobject info = (*env)->CallStaticObjectMethod(env, helper, mid, args);

    fprintf(stderr, "!!!!!!!-> aes: after call staitc obj method\n");
    if (!info) {
        printf("getDecryptedPrice error\n");
        return 0;
    }
    fprintf(stderr, "befor get price encrypt info\n");
    jclass infoCls = getPriceEncryptInfo();
    fprintf(stderr, "befor get price field id \n");
    jfieldID priceId = getPriceFieldId(infoCls);
    fprintf(stderr, "after get price field id \n");
    return (long)(*env)->GetLongField(env, info, priceId);
}

void jprice_init()
{
    JVM_Init();
}

void jprice_destroy()
{
    JVM_Destroy();
}
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
	//"net/url"
	//"time"
)

var decodeKey *C.char

func Init() {
	C.jprice_init()
	ePrice := "RlVGLhEdwomHYLliu4pjMTsqNfh66FL3yb0LbmrOnwQ="
	SetKey("0swvdch0ocmsd0m2viy7c0brrrnhmwpu")
	price := GetDecryptedPrice(ePrice)
	fmt.Println("small cook price: ", price)
}

func Destroy() {
	if unsafe.Pointer(decodeKey) != unsafe.Pointer(nil) {
		C.free(unsafe.Pointer(decodeKey))
	}
	C.jprice_destroy()
}

func SetKey(key string) {
	if unsafe.Pointer(decodeKey) != unsafe.Pointer(nil) {
		C.free(unsafe.Pointer(decodeKey))
	}
	decodeKey = C.CString(key)
}

var mutex sync.Mutex

func GetDecryptedPrice(code string) int {
	//_ = mutex
	mutex.Lock()
	defer mutex.Unlock()

	fmt.Println("code = ", code)
	encryPrice := C.CString(code)
	fmt.Println("get c price = ", encryPrice)
	defer C.free(unsafe.Pointer(encryPrice))

	fmt.Println("aes: this decode key = ", C.GoString(decodeKey))
	return int(C.getDecryptedPrice(encryPrice, decodeKey))
}

// func main() {
// 	ePrice := "RlVGLhEdwomHYLliu4pjMTsqNfh66FL3yb0LbmrOnwQ="
//
// 	now := time.Now()
// 	Init()
// 	fmt.Println("Init use time: ", time.Since(now))
//
// 	SetKey("0swvdch0ocmsd0m2viy7c0brrrnhmwpu")
//
// 	now = time.Now()
// 	price := GetDecryptedPrice(ePrice)
// 	fmt.Println("Decrypted use time: ", time.Since(now))
//
// 	fmt.Println("price: ", price)
//
// 	now = time.Now()
// 	Destroy()
// 	fmt.Println("Destroy use time: ", time.Since(now))
// }

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
    jclass strClass = (*env)->FindClass(env, "java/lang/String");
    jmethodID ctorID = (*env)->GetMethodID(env, strClass, "<init>", "([BLjava/lang/String;)V");
    jbyteArray bytes = (*env)->NewByteArray(env, strlen(pat));
    (*env)->SetByteArrayRegion(env, bytes, 0, strlen(pat), (jbyte*)pat);
    jstring encoding = (*env)->NewStringUTF(env, "utf-8");
    return (jstring)(*env)->NewObject(env, strClass, ctorID, bytes, encoding);
}

static jclass getPriceEncryptHelper()
{
    static jclass helper = NULL;

    if (helper) {
        return helper;
    }

    helper = (*env)->FindClass(env, "PriceEncryptHelper");
    if (!helper) {
        printf("cannot find PriceEncryptHelper\n");
        exit(-1);
    }
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
        return mid;
    }

    mid = (*env)->GetStaticMethodID(env, helper, "getDecryptedPrice", "(Ljava/lang/String;Ljava/lang/String;)LPriceEncryptInfo;");
    if (!mid) {
        printf("cannot find method getDecryptedPriceMethod\n");
        exit(-1);
    }
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
    jclass helper = getPriceEncryptHelper();
    jmethodID mid = getDecryptedPriceMethod(helper);
    jstring encryptedPrice = ch2jstring(ePrice);
    jstring jkey = ch2jstring(key);
    jobject info = (*env)->CallStaticObjectMethod(env, helper, mid, encryptedPrice, jkey);
    if (!info) {
        printf("getDecryptedPrice error\n");
        return 0;
    }
    jclass infoCls = getPriceEncryptInfo();
    jfieldID priceId = getPriceFieldId(infoCls);
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
	//"fmt"
	//"time"
	"unsafe"
)

var decodeKey *C.char

func Init() {
	C.jprice_init()
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

func GetDecryptedPrice(code string) int {
	eprice := C.CString(code)
	defer C.free(unsafe.Pointer(eprice))
	return int(C.getDecryptedPrice(eprice, decodeKey))
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

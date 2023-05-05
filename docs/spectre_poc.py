import struct
import hmac
import hashlib
from hashlib import scrypt as hash

def scrypt(password, salt, N, r, p, dk_len):
    return hash(password, salt=salt, n=N, r=r, p=p, dklen=dk_len, maxmem=(128 * r * (N + p + 2)))

def new_user_key(username, password, algo_v):
  password_bytes = password.encode()
  username_bytes = username.encode()
  purpose = 'com.lyndir.masterpassword'
  purpose_bytes = purpose.encode()

  user_salt = [0] * ( len(purpose_bytes) + 4 + len(username_bytes) )

  uS=0 # start from
  for ctr,sgl in enumerate(purpose_bytes):
    user_salt[ctr] = sgl

  uS += len(purpose_bytes)
  for ctr, sgl in  enumerate(struct.pack('>I', len(username_bytes))):
    user_salt[ctr + uS] = sgl

  uS += 4
  for ctr,sgl in enumerate(username_bytes):
    user_salt[ctr + uS] = sgl

  user_key_data = scrypt(password_bytes, bytes(user_salt), 32768, 8, 2, 64)
  return user_key_data

def new_site_key(user_key_crypto,site_name, key_counter, purpose, key_context):
  site_name_bytes = site_name.encode()
  purpose_bytes = purpose.encode()

  key_context_bytes = None
  key_context_len = 0
  if key_context is not None:
    key_context_bytes = key_context.encode()
    key_context_len = len(key_context_bytes)
  # init
  site_salt = [0] * ( len(purpose_bytes) + 4 + len(site_name_bytes) + 4 +  key_context_len)

  sS=0 # start from
  for ctr,sgl in enumerate(purpose_bytes):
    site_salt[ctr] = sgl
  sS += len(purpose_bytes)

  for ctr, sgl in  enumerate(struct.pack('>I', len(site_name_bytes))):
    site_salt[ctr + sS] = sgl
  sS+=4

  for ctr, sgl in  enumerate(site_name_bytes):
    site_salt[ctr + sS] = sgl
  sS+= len(site_name_bytes)


  for ctr, sgl in  enumerate(struct.pack('>I', key_counter)):
    site_salt[ctr + sS] = sgl
  sS+=4
  key_data = hmac.new(user_key_crypto, msg=bytes(site_salt),digestmod=hashlib.sha256).digest()
  return list(key_data)

templates = {
    "med": ["CvcnoCvc","CvcCvcno"],
    "long":[
        "CvcvnoCvcvCvcv",
        "CvcvCvcvnoCvcv",
        "CvcvCvcvCvcvno",
        "CvccnoCvcvCvcv",
        "CvccCvcvnoCvcv",
        "CvccCvcvCvcvno",
        "CvcvnoCvccCvcv",
        "CvcvCvccnoCvcv",
        "CvcvCvccCvcvno",
        "CvcvnoCvcvCvcc",
        "CvcvCvcvnoCvcc",
        "CvcvCvcvCvccno",
        "CvccnoCvccCvcv",
        "CvccCvccnoCvcv",
        "CvccCvccCvcvno",
        "CvcvnoCvccCvcc",
        "CvcvCvccnoCvcc",
        "CvcvCvccCvccno",
        "CvccnoCvcvCvcc",
        "CvccCvcvnoCvcc",
        "CvccCvcvCvccno"
    ],
}
characters = {
    "V": "AEIOU",
    "C": "BCDFGHJKLMNPQRSTVWXYZ",
    "v": "aeiou",
    "c": "bcdfghjklmnpqrstvwxyz",
    "A": "AEIOUBCDFGHJKLMNPQRSTVWXYZ",
    "a": "AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz",
    "n": "0123456789",
    "o": "@&%?,=[]_:-+*$#!'^~;()/.",
    "x": "AEIOUaeiouBCDFGHJKLMNPQRSTVWXYZbcdfghjklmnpqrstvwxyz0123456789!@#$%^&*()",
    ' ': " "
}

def new_site_result(username, password, site, key_counter=1, key_purpose='com.lyndir.masterpassword', key_type='med'):
  user_key =  new_user_key(username, password,'')
  site_key_bytes = new_site_key(user_key, site, key_counter, key_purpose, None)

  res_templates = templates[key_type]
  res_template = res_templates[site_key_bytes[0] % len(res_templates)]
  pass_res = ""
  for i in range(len(res_template)):
    curr_char = characters[ res_template[i]]
    pass_res += curr_char[ site_key_bytes[i+1] % len(curr_char) ]
  return pass_res

assert new_site_result(username="apho",password="apho", site="fb.com", key_counter=1, key_type="long") == "HargWujtSuya3-"
assert new_site_result(username="apho",password="apho", site="periplus23.com", key_counter=23, key_type="long") == "SujaCupt8:Yovu"
assert new_site_result(username="a",password="a", site="twitter.com", key_counter=1, key_type="med") == "RevXep5+"
assert new_site_result(username="a",password="a", site="twitter.com", key_counter=1, key_type="long") == "RevoGupsWunl3-"

while True:
  uname = input('username -> ')
  passwd = input('password -> ')
  site = input('site -> ' )
  kc = input('counter -> ' )
  kt = input('key type [med|long] -> ' )
  gen_passwd = new_site_result(username=uname, password=passwd, site=site, key_counter=int(kc), key_type=kt)
  print("Generated passwd : " + gen_passwd)


#! /bin/bash


# http://lanlingzi.cn/post/technical/2021/0313_code/
# https://ctftime.org/writeup/16723
# https://github.com/facebookresearch/CrypTen/blob/main/crypten/common/serial.py
#https://huggingface.co/docs/hub/security-pickle

import io
import pickle

class WhiteListUnpickler(pickle.Unpickler):
    def find_class(self, module, name):
        self.check_safe_module(module, name)
        return super().find_class(module, name)
    
    def check_safe_module(module, name):
        # 检查是否在白名单
        if module != '__main__': 
            raise pickle.UnpicklingError("'%s.%s' is forbidden" % (module, name))

def safe_load(file, *, fix_imports=True, encoding="ASCII", errors="strict",
          buffers=None):
    return WhiteListUnpickler(file, fix_imports=fix_imports, buffers=buffers,
                     encoding=encoding, errors=errors).load()

def safe_loads(s, *, fix_imports=True, encoding="ASCII", errors="strict",
           buffers=None):
        file = io.BytesIO(s)
        return WhiteListUnpickler(file, fix_imports=fix_imports, buffers=buffers,
                        encoding=encoding, errors=errors).load()

pickle.load = safe_load
pickle.loads = safe_loads

if __name__ == '__main__':
    bs = b'....' # 需要unpickler的内容
    file = io.BytesIO(bs)
    WhiteListUnpickler(file).load()

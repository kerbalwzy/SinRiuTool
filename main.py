# -*- coding=utf-8 -*-
# @Author: wzy
# @Time: 2021/4/30
"""
1.打开一个wxPython+CEFPython的GUI窗口，并访问配置文件中的URL，进行网页渲染
2.通过Subprocess开启Web服务的后端程序，并记录进程ID
3.浏览器退出时，同时根据后端进程ID关闭所有的后端程序
"""
import platform
import sys

import wx
from cefpython3 import cefpython as cef
# 常量
WINDOWS = (platform.system() == "Windows")
LINUX = (platform.system() == "Linux")
MAC = (platform.system() == "Darwin")
# 运行前检查环境
print("CEF Python {ver}".format(ver=cef.__version__))
print("Python {ver} {arch}".format(ver=platform.python_version(), arch=platform.architecture()[0]))
print("wxPython {ver}".format(ver=wx.version()))
assert cef.__version__ >= "66.0", "要求使用CEF Python v66.0+的版本"
try:
    cef.GetVersion()
except Exception as e:
    raise Exception("{err}\nPython解释器版本在当前操作系统上不支持CEFPython66.0，请重新选择解释器".format(err=e))
if MAC:
    try:
        # noinspection PyUnresolvedReferences
        from AppKit import NSApp
        # Make the content view for the window have a layer.
        # This will make all sub-views have layers. This is
        # necessary to ensure correct layer ordering of all
        # child views and their layers. This fixes Window
        # glitchiness during initial loading on Mac (Issue #371).
        NSApp.windows()[0].contentView().setWantsLayer_(True)
    except ImportError:
        raise Exception("Error: PyObjC package is missing, cannot fix Issue #371\n"
                        "To install PyObjC type: pip install -U pyobjc")
if WINDOWS:
    # noinspection PyUnresolvedReferences, PyArgumentList
    cef.DpiAware.EnableHighDpiSupport()
# 替换异常处理逻辑，保证异常发生时能够结束所有进程
sys.excepthook = cef.ExceptHook

# 配置常量
WIDTH = 1000
HEIGHT = 618

# 全局变量
GCountWindows = 0


def scale_window_size_for_high_dpi(width, height):
    """Scale window size for high DPI devices. This func can be
    called on all operating systems, but scales only for Windows.
    If scaled value is bigger than the work area on the display
    then it will be reduced."""
    if not WINDOWS:
        return width, height
    (_, _, max_width, max_height) = wx.GetClientDisplayRect().Get()
    # noinspection PyUnresolvedReferences
    (width, height) = cef.DpiAware.Scale((width, height))
    if width > max_width:
        width = max_width
    if height > max_height:
        height = max_height
    return width, height


class MainFrame(wx.Frame):

    def __init__(self, parent=None, id=None, title=None, pos=None, size=None, style=None, name=None):
        self.browser = None
        # Must ignore X11 errors like 'BadWindow' and others by
        # installing X11 error handlers. This must be done after
        # wx was initialized.
        if LINUX:
            cef.WindowUtils.InstallX11ErrorHandlers()
        size = size or (WIDTH, HEIGHT)
        size = scale_window_size_for_high_dpi(*size)
        super().__init__(parent, id, title, pos, size, style, name)

        global GCountWindows
        GCountWindows += 1


class UIMain:
    def __init__(self, index_url: str):
        self.index_url = index_url
        self.cef = cef
        pass

    def init_window(self):
        self.cef.Initialize()
        self.cef.CreateBrowserSync(url=self.index_url)
        self.cef.MessageLoop()

    def close_window(self):
        self.cef.Shutdown()

    def run(self):
        self.init_window()
        self.close_window()

        pass


if __name__ == '__main__':
    UIMain(index_url="https://www.baidu.com").run()

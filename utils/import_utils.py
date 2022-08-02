import importlib.util
import os
import sys
from itertools import chain
from types import ModuleType
from typing import Any

from packaging import version

from . import logging

if sys.version_info < (3, 8):
    import importlib_metadata
else:
    import importlib.metadata as importlib_metadata

logger = logging.get_logger(__name__)

_torch_available = importlib.util.find_spec("torch") is not None
if _torch_available:
    try:
        _torch_version = importlib_metadata.version("torch")
        logger.info(f"PyTorch version {_torch_version} available.")
    except importlib_metadata.PackageNotFoundError:
        _torch_available = False

_lightbgm_available = importlib.util.find_spec("lightgbm") is not None
if _lightbgm_available:
    try:
        _lightbgm_version = importlib_metadata.version("lightgbm")
        logger.info(f"LightGBM version {_lightbgm_version} available.")
    except importlib_metadata.PackageNotFoundError:
        _lightbgm_available = False


def is_torch_available():
    return _torch_available


def is_torch_cuda_available():
    if is_torch_available():
        import torch

        return torch.cuda.is_available()
    else:
        return False


def is_torch_bf16_gpu_available():
    if not is_torch_available():
        return False

    import torch
    if version.parse(torch.__version__) < version.parse("1.10"):
        return False

    if torch.cuda.is_available() and torch.version.cuda is not None:
        if torch.cuda.get_device_properties(torch.cuda.current_device()).major < 8:
            return False
        if int(torch.version.cuda.split(".")[0]) < 11:
            return False
        if not hasattr(torch.cuda.amp, "autocast"):
            return False
    else:
        return False

    return True


def is_torch_bf16_cpu_available():
    if not is_torch_available():
        return False

    import torch

    if version.parse(torch.__version__) < version.parse("1.10"):
        return False

    try:
        _ = torch.cpu.amp.autocast
    except AttributeError:
        return False

    return True


def is_torch_tf32_available():
    if not is_torch_available():
        return False

    import torch

    if not torch.cuda.is_available() or torch.version.cuda is None:
        return False
    if torch.cuda.get_device_properties(torch.cuda.current_device()).major < 8:
        return False
    if int(torch.version.cuda.split(".")[0]) < 11:
        return False
    if version.parse(torch.__version__) < version.parse("1.7"):
        return False

    return True


def is_torch_tpu_available(check_device=True):
    if not _torch_available:
        return False
    if importlib.util.find_spec("torch_xla") is not None:
        if check_device:
            try:
                import torch_xla.core.xla_model as xm
                _ = xm.xla_device()
                return True
            except RuntimeError:
                return False
        return True
    return False


def is_lightbgm_available():
    return _lightbgm_available


class _LazyModule(ModuleType):
    def __init__(self, name, module_file, import_structure, module_spec=None, extra_objects=None):
        super().__init__(name)
        self._modules = set(import_structure.keys())
        self._class_to_module = {}
        for key, values in import_structure.items():
            for value in values:
                self._class_to_module[value] = key
        self.__all__ = list(import_structure.keys()) + list(chain(*import_structure.values()))
        self.__file__ = module_file
        self.__spec__ = module_spec
        self.__path__ = [os.path.dirname(module_file)]
        self._objects = {} if extra_objects is None else extra_objects
        self._name = name
        self._import_structure = import_structure

    def __dir__(self):
        result = super().__dir__()
        for attr in self.__all__:
            if attr not in result:
                result.append(attr)
        return result

    def __getattr__(self, name: str) -> Any:
        if name in self._objects:
            return self._objects[name]
        if name in self._modules:
            value = self._get_module(name)
        elif name in self._class_to_module.keys():
            module = self._get_module(self._class_to_module[name])
            value = getattr(module, name)
        else:
            raise AttributeError(f"module {self.__name__} has no attribute {name}")

        setattr(self, name, value)
        return value

    def _get_module(self, module_name: str):
        try:
            return importlib.import_module("." + module_name, self.__name__)
        except Exception as e:
            raise RuntimeError(
                f"Failed to import {self.__name__}.{module_name} because of the following error (look up to see its"
                f" traceback):\n{e}"
            ) from e

    def __reduce__(self):
        return self.__class__, (self._name, self.__file__, self._import_structure)

from setuptools import setup, find_packages

setup(
    name='chess-engine',
    version='0.1',
    packages=find_packages(),
    entry_points={
        'console_scripts': [
            'chess-engine=main:main'
        ]
    },
)
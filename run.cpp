#include "bits/stdc++.h"
using namespace std;

class Base
{
public:
    virtual string name() const;
};

string Base::name() const { return "Base"; }

class Derived : public Base
{
public:
    string name() const override;
};

string Derived::name() const { return "Derived"; }

class MoreDerived : public Derived
{
public:
    string name() const override;
};

string MoreDerived::name() const { return "MoreDerived"; }

int main()
{
    Base base;
    Derived derived;
    MoreDerived morederived;

    cout << base.name() << "  " << derived.name() << "   " << morederived.name() << endl;
}

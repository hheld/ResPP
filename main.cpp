#include <iostream>
#include <string>

extern "C"
{
#include <b64.h>
}

using namespace std;

int main(int argc, char *argv[])
{
    (void) argc;
    (void) argv;

    std::string s("Hello, world!");

    unsigned long numEncodedBytes;

    auto s_enc = b64_encode(s.c_str(), s.size(), &numEncodedBytes);
    auto s_enc_dec = b64_decode(s_enc, numEncodedBytes, nullptr);

    cout << "original: " << s << "\n";
    cout << "encoded: " << s_enc << "\n";
    cout << "decoded: " << s_enc_dec << "\n";

    free(s_enc);
    free(s_enc_dec);

    return 0;
}

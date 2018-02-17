#define CATCH_CONFIG_MAIN

#include /* github.com/catchorg/Catch2/include/ */ "catch.hpp"
#include/*https://github.com/nothings/stb*/"stb.h"

TEST_CASE("STB is included" ) {
    REQUIRE( stb_min(10, 11) == 10 );
}

#include "fs.h"
#import <Foundation/Foundation.h>

void
get_dir(int directory, int domain, char* buf, int maxlen)
{
	NSURL* url = [[NSFileManager defaultManager] URLForDirectory:directory inDomain:domain appropriateForURL:nil create:NO error:nil];
	[url getFileSystemRepresentation:buf maxLength:maxlen];
}


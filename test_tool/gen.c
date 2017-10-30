#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

#include <bliss_b_keys.h>

void gen_private_key(entropy_t *entropy, bliss_param_t *params) {
	int i;

	bliss_private_key_t private_key;

	/* Generate private key */
	bliss_b_private_key_gen(&private_key, params->kind, entropy);

	/* Output private key */
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s1[i]);
	printf("\n");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.s2[i]);
	printf("\n");
	for(i = 0; i < params->n; i++)
		printf("%d ",private_key.a[i]);
	printf("\n");

	bliss_b_private_key_delete(&private_key);

}

void usage(char *argv[]) {
	fprintf(stderr,"Usage: %s subcommand [options]\n",argv[0]);
	fprintf(stderr,"  subcommands\n");
	fprintf(stderr,"    keygen\n");
	fprintf(stderr,"  options\n");
	fprintf(stderr,"    -k kind  0/1/2/3/4 which version of BLISS?\n");
	exit(1);
}

const int GEN_KEY = 0;

int main(int argc, char *argv[]) {
	int i,c;
	int kind = 0;
	int subcommand = -1;
	entropy_t entropy;
	uint8_t seed[64];
	bliss_param_t params;

	/* Parse command line */
	if(argc < 2) {
		usage(argv);
	} else {
		if(strcmp(argv[1],"keygen") == 0) {
			subcommand = GEN_KEY;
			fprintf(stderr,"Generating private key...\n");
		} else {
			subcommand = -1;
		}
	}

	while((c = getopt(argc-1,argv+1,"k:")) != -1) {
		switch(c) {
			case 'k':
				kind = atoi(optarg);
				if(kind < 0 || kind > 4)
					usage(argv);
				fprintf(stderr,"Version set to %d\n",kind);
				break;
		}
	}

	/* Initialize entropy */
	for(i = 0; i < 64; i++)
		seed[i] = i%8;
	entropy_init(&entropy,seed);

	/* Initialize parameter */
	bliss_params_init(&params, kind);

	if(subcommand == GEN_KEY) {
		gen_private_key(&entropy,&params);
	} else {
		usage(argv);
	}

	return 0;
}
